// Copyright (c) 2023 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package server

import (
	"time"

	validator "github.com/AccelByte/justice-input-validation-go"
	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
)

const (
	WeaponTypeSword = "SWORD"
	WeaponTypeGun   = "GUN"

	MaxSizeEventBannerInKB = 100
)

type CustomGameRecord struct {
	LocationID     string `json:"locationId" valid:"required~locationId cannot be empty"`
	Name           string `json:"name" valid:"required~name cannot be empty"`
	TotalResources int    `json:"totalResources" valid:"required~totalResources cannot be empty"`
	TotalEnemy     int    `json:"totalEnemy" valid:"required~totalEnemy cannot be empty"`
}

func (c *CustomGameRecord) Validate() error {
	_, err := validator.ValidateStruct(c)
	if err != nil {
		return err
	}

	return nil
}

type CustomPlayerRecord struct {
	UserID              string `json:"userId" valid:"required~user ID cannot be empty"`
	FavouriteWeaponType string `json:"favouriteWeaponType" valid:"required~favourite weapon type cannot be empty"`
	FavouriteWeapon     string `json:"favouriteWeapon" valid:"required~favourite weapon cannot be empty"`
}

func (c *CustomPlayerRecord) Validate() error {
	_, err := govalidator.ValidateStruct(c)
	if err != nil {
		return err
	}
	if c.FavouriteWeaponType != WeaponTypeSword && c.FavouriteWeaponType != WeaponTypeGun {
		return errors.New("invalid weapon type")
	}

	return nil
}

type DailyMessage struct {
	Message     string    `json:"message"`
	Title       string    `json:"title"`
	AvailableOn time.Time `json:"availableOn"`
}

type PlayerActivity struct {
	UserID   string `json:"userId" valid:"required~user ID cannot be empty"`
	Activity string `json:"activity" valid:"required~activity cannot be empty"`
}

func (c *PlayerActivity) Validate() error {
	_, err := govalidator.ValidateStruct(c)
	if err != nil {
		return err
	}

	return nil
}
