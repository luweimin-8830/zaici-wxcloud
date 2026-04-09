package service

import (
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
)

// MediaService 媒体服务层
type MediaService struct {
	pictureDao *dao.PictureListDao
}

// NewMediaService 创建媒体服务实例
func NewMediaService() *MediaService {
	return &MediaService{
		pictureDao: &dao.PictureListDao{},
	}
}

// CreatePicture 创建图片记录
func (s *MediaService) CreatePicture(userPicURL, picHash, traceID string) (*model.PictureList, error) {
	return s.CreatePictureWithStatus(userPicURL, picHash, traceID, 0)
}

// CreatePictureWithStatus 创建图片记录并指定审核状态
func (s *MediaService) CreatePictureWithStatus(userPicURL, picHash, traceID string, status int) (*model.PictureList, error) {
	picture := &model.PictureList{
		UserPicURL:     userPicURL,
		PicHash:        picHash,
		TraceID:        traceID,
		SecCheckStatus: status,
	}
	err := s.pictureDao.Create(picture)
	if err != nil {
		return nil, err
	}
	return picture, nil
}

// GetPictureByID 根据ID获取图片
func (s *MediaService) GetPictureByID(id uint) (*model.PictureList, error) {
	return s.pictureDao.GetByID(id)
}

// GetPictureByHash 根据图片哈希获取图片
func (s *MediaService) GetPictureByHash(picHash string) (*model.PictureList, error) {
	return s.pictureDao.GetByPicHash(picHash)
}

// UpdatePicture 更新图片记录
func (s *MediaService) UpdatePicture(picture *model.PictureList) error {
	return s.pictureDao.Update(picture)
}

// UpdatePictureSecCheckStatus 更新图片审核状态
func (s *MediaService) UpdatePictureSecCheckStatus(id uint, status int) error {
	return s.pictureDao.UpdateSecCheckStatus(id, status)
}

// DeletePicture 删除图片记录
func (s *MediaService) DeletePicture(id uint) error {
	return s.pictureDao.Delete(id)
}

// ListPictures 获取图片列表
func (s *MediaService) ListPictures(limit, offset int) ([]model.PictureList, error) {
	return s.pictureDao.List(limit, offset)
}

// CountPictures 获取图片总数
func (s *MediaService) CountPictures() (int64, error) {
	return s.pictureDao.Count()
}
