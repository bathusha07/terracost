package terraform_test

import (
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cycloidio/cost-estimation/mock"
	"github.com/cycloidio/cost-estimation/query"
	"github.com/cycloidio/cost-estimation/terraform"
)

func TestPlan_ExtractPlannedQueries(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		provider := mock.NewTerraformProvider(ctrl)

		plan := terraform.NewPlan(terraform.ProviderInitializer{
			MatchNames: []string{"aws-test"},
			Provider: func(_ terraform.ProviderConfig) (terraform.Provider, error) {
				return provider, nil
			},
		})

		f, err := os.Open("../testdata/terraform-plan.json")
		require.NoError(t, err)
		defer f.Close()

		err = plan.Read(f)
		require.NoError(t, err)

		provider.EXPECT().ResourceComponents(gomock.Any()).DoAndReturn(func(res terraform.Resource) ([]query.Component, error) {
			assert.Equal(t, "aws_instance.example", res.Address)
			assert.Equal(t, "t2.xlarge", res.Values["instance_type"])
			return []query.Component{}, nil
		})

		queries, err := plan.ExtractPlannedQueries()
		require.NoError(t, err)
		require.Len(t, queries, 1)
	})

	t.Run("BadProvider", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		plan := terraform.NewPlan(terraform.ProviderInitializer{
			MatchNames: []string{"aws-test"},
			Provider: func(_ terraform.ProviderConfig) (terraform.Provider, error) {
				return nil, errors.New("bad provider")
			},
		})

		f, err := os.Open("../testdata/terraform-plan.json")
		require.NoError(t, err)
		defer f.Close()

		err = plan.Read(f)
		require.NoError(t, err)

		_, err = plan.ExtractPlannedQueries()
		assert.Error(t, err)
	})

	t.Run("FailedResourceComponents", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		provider := mock.NewTerraformProvider(ctrl)

		plan := terraform.NewPlan(terraform.ProviderInitializer{
			MatchNames: []string{"aws-test"},
			Provider: func(_ terraform.ProviderConfig) (terraform.Provider, error) {
				return provider, nil
			},
		})

		f, err := os.Open("../testdata/terraform-plan.json")
		require.NoError(t, err)
		defer f.Close()

		err = plan.Read(f)
		require.NoError(t, err)

		provider.EXPECT().ResourceComponents(gomock.Any()).DoAndReturn(func(res terraform.Resource) ([]query.Component, error) {
			return nil, errors.New("ResourceComponents fail")
		})

		queries, err := plan.ExtractPlannedQueries()
		require.NoError(t, err)
		require.Len(t, queries, 0)
	})
}

func TestPlan_ExtractPriorQueries(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		provider := mock.NewTerraformProvider(ctrl)

		plan := terraform.NewPlan(terraform.ProviderInitializer{
			MatchNames: []string{"aws-test"},
			Provider: func(_ terraform.ProviderConfig) (terraform.Provider, error) {
				return provider, nil
			},
		})

		f, err := os.Open("../testdata/terraform-plan.json")
		require.NoError(t, err)
		defer f.Close()

		err = plan.Read(f)
		require.NoError(t, err)

		provider.EXPECT().ResourceComponents(gomock.Any()).DoAndReturn(func(res terraform.Resource) ([]query.Component, error) {
			assert.Equal(t, "aws_instance.example", res.Address)
			assert.Equal(t, "t2.micro", res.Values["instance_type"])
			return []query.Component{}, nil
		})

		queries, err := plan.ExtractPriorQueries()
		require.NoError(t, err)
		require.Len(t, queries, 1)
	})

	t.Run("BadProvider", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		plan := terraform.NewPlan(terraform.ProviderInitializer{
			MatchNames: []string{"aws-test"},
			Provider: func(_ terraform.ProviderConfig) (terraform.Provider, error) {
				return nil, errors.New("bad provider")
			},
		})

		f, err := os.Open("../testdata/terraform-plan.json")
		require.NoError(t, err)
		defer f.Close()

		err = plan.Read(f)
		require.NoError(t, err)

		_, err = plan.ExtractPriorQueries()
		assert.Error(t, err)
	})

	t.Run("FailedResourceComponents", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		provider := mock.NewTerraformProvider(ctrl)

		plan := terraform.NewPlan(terraform.ProviderInitializer{
			MatchNames: []string{"aws-test"},
			Provider: func(_ terraform.ProviderConfig) (terraform.Provider, error) {
				return provider, nil
			},
		})

		f, err := os.Open("../testdata/terraform-plan.json")
		require.NoError(t, err)
		defer f.Close()

		err = plan.Read(f)
		require.NoError(t, err)

		provider.EXPECT().ResourceComponents(gomock.Any()).DoAndReturn(func(res terraform.Resource) ([]query.Component, error) {
			return nil, errors.New("ResourceComponents fail")
		})

		queries, err := plan.ExtractPriorQueries()
		require.NoError(t, err)
		require.Len(t, queries, 0)
	})
}