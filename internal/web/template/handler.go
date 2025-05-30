package template

import (
	"errors"
	"log"

	"gitee.com/flycash/notification-platform/internal/domain"
	"gitee.com/flycash/notification-platform/internal/errs"
	templatesvc "gitee.com/flycash/notification-platform/internal/service/template/manage"
	"github.com/ecodeclub/ekit/slice"
	"github.com/ecodeclub/ginx"

	"github.com/gin-gonic/gin"
)

var _ ginx.Handler = &Handler{}

type Handler struct {
	svc templatesvc.ChannelTemplateService
}

func NewHandler(svc templatesvc.ChannelTemplateService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) PrivateRoutes(_ *gin.Engine) {
}

func (h *Handler) PublicRoutes(server *gin.Engine) {
	g := server.Group("/templates")
	g.POST("/list", ginx.B[ListTemplatesReq](h.ListTemplates))
	g.POST("/create", ginx.B[CreateTemplateReq](h.CreateTemplate))
	g.POST("/update", ginx.B[UpdateTemplateReq](h.UpdateTemplate))
	g.POST("/publish", ginx.B[PublishTemplateReq](h.PublishTemplate))

	j := g.Group("/versions")
	j.POST("/fork", ginx.B[ForkVersionReq](h.ForkVersion))
	j.POST("/update", ginx.B[UpdateVersionReq](h.UpdateVersion))
	j.POST("/review/internal", ginx.B[SubmitForInternalReviewReq](h.SubmitForInternalReview))
}

// ListTemplates 获取所有模版
func (h *Handler) ListTemplates(ctx *ginx.Context, req ListTemplatesReq) (ginx.Result, error) {
	templates, err := h.svc.GetTemplatesByOwner(ctx.Request.Context(), req.OwnerID, domain.OwnerType(req.OwnerType))
	if err != nil {
		return systemErrorResult, err
	}
	return ginx.Result{
		Data: ListTemplatesResp{
			Templates: slice.Map(templates, func(_ int, src domain.ChannelTemplate) ChannelTemplate {
				return h.toTemplateVO(src)
			}),
		},
	}, nil
}

func (h *Handler) toTemplateVO(src domain.ChannelTemplate) ChannelTemplate {
	return ChannelTemplate{
		ID:              src.ID,
		OwnerID:         src.OwnerID,
		OwnerType:       src.OwnerType.String(),
		Name:            src.Name,
		Description:     src.Description,
		Channel:         src.Channel.String(),
		BusinessType:    src.BusinessType.ToInt64(),
		ActiveVersionID: src.ActiveVersionID,
		Ctime:           src.Ctime,
		Utime:           src.Utime,
		Versions: slice.Map(src.Versions, func(_ int, src domain.ChannelTemplateVersion) ChannelTemplateVersion {
			return h.toVersionVO(src)
		}),
	}
}

func (h *Handler) toVersionVO(src domain.ChannelTemplateVersion) ChannelTemplateVersion {
	return ChannelTemplateVersion{
		ID:                       src.ID,
		ChannelTemplateID:        src.ChannelTemplateID,
		Name:                     src.Name,
		Signature:                src.Signature,
		Content:                  src.Content,
		Remark:                   src.Remark,
		AuditID:                  src.AuditID,
		AuditorID:                src.AuditorID,
		AuditTime:                src.AuditTime,
		AuditStatus:              src.AuditStatus.String(),
		RejectReason:             src.RejectReason,
		LastReviewSubmissionTime: src.LastReviewSubmissionTime,
		Ctime:                    src.Ctime,
		Utime:                    src.Utime,
		Providers: slice.Map(src.Providers, func(_ int, src domain.ChannelTemplateProvider) ChannelTemplateProvider {
			return h.toProviderVO(src)
		}),
	}
}

func (h *Handler) toProviderVO(src domain.ChannelTemplateProvider) ChannelTemplateProvider {
	return ChannelTemplateProvider{
		ID:                       src.ID,
		TemplateID:               src.TemplateID,
		TemplateVersionID:        src.TemplateVersionID,
		ProviderID:               src.ProviderID,
		ProviderName:             src.ProviderName,
		ProviderChannel:          src.ProviderChannel.String(),
		RequestID:                src.RequestID,
		ProviderTemplateID:       src.ProviderTemplateID,
		AuditStatus:              src.AuditStatus.String(),
		RejectReason:             src.RejectReason,
		LastReviewSubmissionTime: src.LastReviewSubmissionTime,
		Ctime:                    src.Ctime,
		Utime:                    src.Utime,
	}
}

// CreateTemplate 创建模板
func (h *Handler) CreateTemplate(ctx *ginx.Context, req CreateTemplateReq) (ginx.Result, error) {
	template := domain.ChannelTemplate{
		OwnerID:      req.OwnerID,
		OwnerType:    domain.OwnerType(req.OwnerType),
		Name:         req.Name,
		Description:  req.Description,
		Channel:      domain.Channel(req.Channel),
		BusinessType: domain.BusinessType(req.BusinessType),
	}

	createdTemplate, err := h.svc.CreateTemplate(ctx.Request.Context(), template)
	if err != nil {
		return systemErrorResult, err
	}

	return ginx.Result{
		Data: CreateTemplateResp{
			Template: h.toTemplateVO(createdTemplate),
		},
	}, nil
}

// UpdateTemplate 更新模板基础信息
func (h *Handler) UpdateTemplate(ctx *ginx.Context, req UpdateTemplateReq) (ginx.Result, error) {
	template := domain.ChannelTemplate{
		ID:           req.TemplateID,
		Name:         req.Name,
		Description:  req.Description,
		BusinessType: domain.BusinessType(req.BusinessType),
	}

	if err := h.svc.UpdateTemplate(ctx.Request.Context(), template); err != nil {

		log.Printf("err = %#v\n", err)
		return systemErrorResult, err
	}

	return ginx.Result{
		Msg: "OK",
	}, nil
}

// PublishTemplate 发布模板
func (h *Handler) PublishTemplate(ctx *ginx.Context, req PublishTemplateReq) (ginx.Result, error) {
	if err := h.svc.PublishTemplate(ctx.Request.Context(), req.TemplateID, req.VersionID); err != nil {
		return systemErrorResult, err
	}

	return ginx.Result{
		Msg: "OK",
	}, nil
}

// ForkVersion 拷贝模版版本
func (h *Handler) ForkVersion(ctx *ginx.Context, req ForkVersionReq) (ginx.Result, error) {
	version, err := h.svc.ForkVersion(ctx.Request.Context(), req.VersionID)
	if err != nil {
		return systemErrorResult, err
	}
	return ginx.Result{
		Data: ForkVersionResp{
			TemplateVersion: h.toVersionVO(version),
		},
	}, nil
}

// UpdateVersion 更新模板版本
func (h *Handler) UpdateVersion(ctx *ginx.Context, req UpdateVersionReq) (ginx.Result, error) {
	version := domain.ChannelTemplateVersion{
		ID:        req.VersionID,
		Name:      req.Name,
		Signature: req.Signature,
		Content:   req.Content,
		Remark:    req.Remark,
	}

	if err := h.svc.UpdateVersion(ctx.Request.Context(), version); err != nil {
		if !errors.Is(err, errs.ErrTemplateVersionNotFound) {
			return systemErrorResult, err
		}
	}

	return ginx.Result{
		Msg: "OK",
	}, nil
}

// SubmitForInternalReview 提交内部审核
func (h *Handler) SubmitForInternalReview(ctx *ginx.Context, req SubmitForInternalReviewReq) (ginx.Result, error) {
	if err := h.svc.SubmitForInternalReview(ctx.Request.Context(), req.VersionID); err != nil {
		return systemErrorResult, err
	}

	return ginx.Result{
		Msg: "OK",
	}, nil
}
