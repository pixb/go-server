package store

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"

	storepb "github.com/pixb/go-server/proto/gen/store"
)

type InstanceSetting struct {
	Name        string
	Value       string
	Description string
}

type FindInstanceSetting struct {
	Name string
}

type DeleteInstanceSetting struct {
	Name string
}

func (s *Store) UpsertInstanceSetting(ctx context.Context, upsert *storepb.InstanceSetting) (*storepb.InstanceSetting, error) {
	instanceSettingRaw := &InstanceSetting{
		Name: upsert.Key.String(),
	}
	var valueBytes []byte
	var err error
	if upsert.Key == storepb.InstanceSettingKey_BASIC {
		valueBytes, err = protojson.Marshal(upsert.GetBasicSetting())
	} else {
		return nil, errors.Errorf("unsupported instance setting key: %v", upsert.Key)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal instance setting value")
	}
	valueString := string(valueBytes)
	instanceSettingRaw.Value = valueString
	instanceSettingRaw, err = s.driver.UpsertInstanceSetting(ctx, instanceSettingRaw)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to upsert instance setting")
	}
	instanceSetting, err := convertInstanceSettingFromRaw(instanceSettingRaw)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert instance setting")
	}
	s.instanceSettingCache.Set(ctx, instanceSetting.Key.String(), instanceSetting)
	return instanceSetting, nil
}

func (s *Store) ListInstanceSettings(ctx context.Context, find *FindInstanceSetting) ([]*storepb.InstanceSetting, error) {
	list, err := s.driver.ListInstanceSettings(ctx, find)
	if err != nil {
		return nil, err
	}

	instanceSettings := []*storepb.InstanceSetting{}
	for _, instanceSettingRaw := range list {
		instanceSetting, err := convertInstanceSettingFromRaw(instanceSettingRaw)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to convert instance setting")
		}
		if instanceSetting == nil {
			continue
		}
		s.instanceSettingCache.Set(ctx, instanceSetting.Key.String(), instanceSetting)
		instanceSettings = append(instanceSettings, instanceSetting)
	}
	return instanceSettings, nil
}

func (s *Store) GetInstanceSetting(ctx context.Context, find *FindInstanceSetting) (*storepb.InstanceSetting, error) {
	if cache, ok := s.instanceSettingCache.Get(ctx, find.Name); ok {
		instanceSetting, ok := cache.(*storepb.InstanceSetting)
		if ok {
			return instanceSetting, nil
		}
	}

	list, err := s.ListInstanceSettings(ctx, find)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	if len(list) > 1 {
		return nil, errors.Errorf("found multiple instance settings with key %s", find.Name)
	}
	return list[0], nil
}

func (s *Store) GetInstanceBasicSetting(ctx context.Context) (*storepb.InstanceBasicSetting, error) {
	instanceSetting, err := s.GetInstanceSetting(ctx, &FindInstanceSetting{
		Name: storepb.InstanceSettingKey_BASIC.String(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get instance basic setting")
	}

	instanceBasicSetting := &storepb.InstanceBasicSetting{}
	if instanceSetting != nil {
		instanceBasicSetting = instanceSetting.GetBasicSetting()
	}
	s.instanceSettingCache.Set(ctx, storepb.InstanceSettingKey_BASIC.String(), &storepb.InstanceSetting{
		Key:   storepb.InstanceSettingKey_BASIC,
		Value: &storepb.InstanceSetting_BasicSetting{BasicSetting: instanceBasicSetting},
	})
	return instanceBasicSetting, nil
}

func convertInstanceSettingFromRaw(instanceSettingRaw *InstanceSetting) (*storepb.InstanceSetting, error) {
	instanceSetting := &storepb.InstanceSetting{
		Key: storepb.InstanceSettingKey(storepb.InstanceSettingKey_value[instanceSettingRaw.Name]),
	}
	switch instanceSettingRaw.Name {
	case storepb.InstanceSettingKey_BASIC.String():
		basicSetting := &storepb.InstanceBasicSetting{}
		if err := protojsonUnmarshaler.Unmarshal([]byte(instanceSettingRaw.Value), basicSetting); err != nil {
			return nil, err
		}
		instanceSetting.Value = &storepb.InstanceSetting_BasicSetting{BasicSetting: basicSetting}
	default:
		// Skip unsupported instance setting key.
		return nil, nil
	}
	return instanceSetting, nil
}
