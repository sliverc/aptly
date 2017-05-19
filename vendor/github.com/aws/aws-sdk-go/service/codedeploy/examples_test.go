// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.

package codedeploy_test

import (
	"bytes"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codedeploy"
)

var _ time.Duration
var _ bytes.Buffer

func ExampleCodeDeploy_AddTagsToOnPremisesInstances() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.AddTagsToOnPremisesInstancesInput{
		InstanceNames: []*string{ // Required
			aws.String("InstanceName"), // Required
			// More values...
		},
		Tags: []*codedeploy.Tag{ // Required
			{ // Required
				Key:   aws.String("Key"),
				Value: aws.String("Value"),
			},
			// More values...
		},
	}
	resp, err := svc.AddTagsToOnPremisesInstances(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_BatchGetApplicationRevisions() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.BatchGetApplicationRevisionsInput{
		ApplicationName: aws.String("ApplicationName"), // Required
		Revisions: []*codedeploy.RevisionLocation{ // Required
			{ // Required
				GitHubLocation: &codedeploy.GitHubLocation{
					CommitId:   aws.String("CommitId"),
					Repository: aws.String("Repository"),
				},
				RevisionType: aws.String("RevisionLocationType"),
				S3Location: &codedeploy.S3Location{
					Bucket:     aws.String("S3Bucket"),
					BundleType: aws.String("BundleType"),
					ETag:       aws.String("ETag"),
					Key:        aws.String("S3Key"),
					Version:    aws.String("VersionId"),
				},
			},
			// More values...
		},
	}
	resp, err := svc.BatchGetApplicationRevisions(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_BatchGetApplications() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.BatchGetApplicationsInput{
		ApplicationNames: []*string{
			aws.String("ApplicationName"), // Required
			// More values...
		},
	}
	resp, err := svc.BatchGetApplications(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_BatchGetDeploymentGroups() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.BatchGetDeploymentGroupsInput{
		ApplicationName: aws.String("ApplicationName"), // Required
		DeploymentGroupNames: []*string{ // Required
			aws.String("DeploymentGroupName"), // Required
			// More values...
		},
	}
	resp, err := svc.BatchGetDeploymentGroups(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_BatchGetDeploymentInstances() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.BatchGetDeploymentInstancesInput{
		DeploymentId: aws.String("DeploymentId"), // Required
		InstanceIds: []*string{ // Required
			aws.String("InstanceId"), // Required
			// More values...
		},
	}
	resp, err := svc.BatchGetDeploymentInstances(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_BatchGetDeployments() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.BatchGetDeploymentsInput{
		DeploymentIds: []*string{
			aws.String("DeploymentId"), // Required
			// More values...
		},
	}
	resp, err := svc.BatchGetDeployments(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_BatchGetOnPremisesInstances() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.BatchGetOnPremisesInstancesInput{
		InstanceNames: []*string{
			aws.String("InstanceName"), // Required
			// More values...
		},
	}
	resp, err := svc.BatchGetOnPremisesInstances(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_ContinueDeployment() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.ContinueDeploymentInput{
		DeploymentId: aws.String("DeploymentId"),
	}
	resp, err := svc.ContinueDeployment(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_CreateApplication() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.CreateApplicationInput{
		ApplicationName: aws.String("ApplicationName"), // Required
	}
	resp, err := svc.CreateApplication(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_CreateDeployment() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.CreateDeploymentInput{
		ApplicationName: aws.String("ApplicationName"), // Required
		AutoRollbackConfiguration: &codedeploy.AutoRollbackConfiguration{
			Enabled: aws.Bool(true),
			Events: []*string{
				aws.String("AutoRollbackEvent"), // Required
				// More values...
			},
		},
		DeploymentConfigName:          aws.String("DeploymentConfigName"),
		DeploymentGroupName:           aws.String("DeploymentGroupName"),
		Description:                   aws.String("Description"),
		IgnoreApplicationStopFailures: aws.Bool(true),
		Revision: &codedeploy.RevisionLocation{
			GitHubLocation: &codedeploy.GitHubLocation{
				CommitId:   aws.String("CommitId"),
				Repository: aws.String("Repository"),
			},
			RevisionType: aws.String("RevisionLocationType"),
			S3Location: &codedeploy.S3Location{
				Bucket:     aws.String("S3Bucket"),
				BundleType: aws.String("BundleType"),
				ETag:       aws.String("ETag"),
				Key:        aws.String("S3Key"),
				Version:    aws.String("VersionId"),
			},
		},
		TargetInstances: &codedeploy.TargetInstances{
			AutoScalingGroups: []*string{
				aws.String("AutoScalingGroupName"), // Required
				// More values...
			},
			TagFilters: []*codedeploy.EC2TagFilter{
				{ // Required
					Key:   aws.String("Key"),
					Type:  aws.String("EC2TagFilterType"),
					Value: aws.String("Value"),
				},
				// More values...
			},
		},
		UpdateOutdatedInstancesOnly: aws.Bool(true),
	}
	resp, err := svc.CreateDeployment(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_CreateDeploymentConfig() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.CreateDeploymentConfigInput{
		DeploymentConfigName: aws.String("DeploymentConfigName"), // Required
		MinimumHealthyHosts: &codedeploy.MinimumHealthyHosts{
			Type:  aws.String("MinimumHealthyHostsType"),
			Value: aws.Int64(1),
		},
	}
	resp, err := svc.CreateDeploymentConfig(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_CreateDeploymentGroup() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.CreateDeploymentGroupInput{
		ApplicationName:     aws.String("ApplicationName"),     // Required
		DeploymentGroupName: aws.String("DeploymentGroupName"), // Required
		ServiceRoleArn:      aws.String("Role"),                // Required
		AlarmConfiguration: &codedeploy.AlarmConfiguration{
			Alarms: []*codedeploy.Alarm{
				{ // Required
					Name: aws.String("AlarmName"),
				},
				// More values...
			},
			Enabled:                aws.Bool(true),
			IgnorePollAlarmFailure: aws.Bool(true),
		},
		AutoRollbackConfiguration: &codedeploy.AutoRollbackConfiguration{
			Enabled: aws.Bool(true),
			Events: []*string{
				aws.String("AutoRollbackEvent"), // Required
				// More values...
			},
		},
		AutoScalingGroups: []*string{
			aws.String("AutoScalingGroupName"), // Required
			// More values...
		},
		BlueGreenDeploymentConfiguration: &codedeploy.BlueGreenDeploymentConfiguration{
			DeploymentReadyOption: &codedeploy.DeploymentReadyOption{
				ActionOnTimeout:   aws.String("DeploymentReadyAction"),
				WaitTimeInMinutes: aws.Int64(1),
			},
			GreenFleetProvisioningOption: &codedeploy.GreenFleetProvisioningOption{
				Action: aws.String("GreenFleetProvisioningAction"),
			},
			TerminateBlueInstancesOnDeploymentSuccess: &codedeploy.BlueInstanceTerminationOption{
				Action: aws.String("InstanceAction"),
				TerminationWaitTimeInMinutes: aws.Int64(1),
			},
		},
		DeploymentConfigName: aws.String("DeploymentConfigName"),
		DeploymentStyle: &codedeploy.DeploymentStyle{
			DeploymentOption: aws.String("DeploymentOption"),
			DeploymentType:   aws.String("DeploymentType"),
		},
		Ec2TagFilters: []*codedeploy.EC2TagFilter{
			{ // Required
				Key:   aws.String("Key"),
				Type:  aws.String("EC2TagFilterType"),
				Value: aws.String("Value"),
			},
			// More values...
		},
		LoadBalancerInfo: &codedeploy.LoadBalancerInfo{
			ElbInfoList: []*codedeploy.ELBInfo{
				{ // Required
					Name: aws.String("ELBName"),
				},
				// More values...
			},
		},
		OnPremisesInstanceTagFilters: []*codedeploy.TagFilter{
			{ // Required
				Key:   aws.String("Key"),
				Type:  aws.String("TagFilterType"),
				Value: aws.String("Value"),
			},
			// More values...
		},
		TriggerConfigurations: []*codedeploy.TriggerConfig{
			{ // Required
				TriggerEvents: []*string{
					aws.String("TriggerEventType"), // Required
					// More values...
				},
				TriggerName:      aws.String("TriggerName"),
				TriggerTargetArn: aws.String("TriggerTargetArn"),
			},
			// More values...
		},
	}
	resp, err := svc.CreateDeploymentGroup(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_DeleteApplication() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.DeleteApplicationInput{
		ApplicationName: aws.String("ApplicationName"), // Required
	}
	resp, err := svc.DeleteApplication(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_DeleteDeploymentConfig() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.DeleteDeploymentConfigInput{
		DeploymentConfigName: aws.String("DeploymentConfigName"), // Required
	}
	resp, err := svc.DeleteDeploymentConfig(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_DeleteDeploymentGroup() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.DeleteDeploymentGroupInput{
		ApplicationName:     aws.String("ApplicationName"),     // Required
		DeploymentGroupName: aws.String("DeploymentGroupName"), // Required
	}
	resp, err := svc.DeleteDeploymentGroup(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_DeregisterOnPremisesInstance() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.DeregisterOnPremisesInstanceInput{
		InstanceName: aws.String("InstanceName"), // Required
	}
	resp, err := svc.DeregisterOnPremisesInstance(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_GetApplication() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.GetApplicationInput{
		ApplicationName: aws.String("ApplicationName"), // Required
	}
	resp, err := svc.GetApplication(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_GetApplicationRevision() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.GetApplicationRevisionInput{
		ApplicationName: aws.String("ApplicationName"), // Required
		Revision: &codedeploy.RevisionLocation{ // Required
			GitHubLocation: &codedeploy.GitHubLocation{
				CommitId:   aws.String("CommitId"),
				Repository: aws.String("Repository"),
			},
			RevisionType: aws.String("RevisionLocationType"),
			S3Location: &codedeploy.S3Location{
				Bucket:     aws.String("S3Bucket"),
				BundleType: aws.String("BundleType"),
				ETag:       aws.String("ETag"),
				Key:        aws.String("S3Key"),
				Version:    aws.String("VersionId"),
			},
		},
	}
	resp, err := svc.GetApplicationRevision(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_GetDeployment() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.GetDeploymentInput{
		DeploymentId: aws.String("DeploymentId"), // Required
	}
	resp, err := svc.GetDeployment(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_GetDeploymentConfig() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.GetDeploymentConfigInput{
		DeploymentConfigName: aws.String("DeploymentConfigName"), // Required
	}
	resp, err := svc.GetDeploymentConfig(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_GetDeploymentGroup() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.GetDeploymentGroupInput{
		ApplicationName:     aws.String("ApplicationName"),     // Required
		DeploymentGroupName: aws.String("DeploymentGroupName"), // Required
	}
	resp, err := svc.GetDeploymentGroup(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_GetDeploymentInstance() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.GetDeploymentInstanceInput{
		DeploymentId: aws.String("DeploymentId"), // Required
		InstanceId:   aws.String("InstanceId"),   // Required
	}
	resp, err := svc.GetDeploymentInstance(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_GetOnPremisesInstance() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.GetOnPremisesInstanceInput{
		InstanceName: aws.String("InstanceName"), // Required
	}
	resp, err := svc.GetOnPremisesInstance(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_ListApplicationRevisions() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.ListApplicationRevisionsInput{
		ApplicationName: aws.String("ApplicationName"), // Required
		Deployed:        aws.String("ListStateFilterAction"),
		NextToken:       aws.String("NextToken"),
		S3Bucket:        aws.String("S3Bucket"),
		S3KeyPrefix:     aws.String("S3Key"),
		SortBy:          aws.String("ApplicationRevisionSortBy"),
		SortOrder:       aws.String("SortOrder"),
	}
	resp, err := svc.ListApplicationRevisions(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_ListApplications() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.ListApplicationsInput{
		NextToken: aws.String("NextToken"),
	}
	resp, err := svc.ListApplications(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_ListDeploymentConfigs() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.ListDeploymentConfigsInput{
		NextToken: aws.String("NextToken"),
	}
	resp, err := svc.ListDeploymentConfigs(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_ListDeploymentGroups() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.ListDeploymentGroupsInput{
		ApplicationName: aws.String("ApplicationName"), // Required
		NextToken:       aws.String("NextToken"),
	}
	resp, err := svc.ListDeploymentGroups(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_ListDeploymentInstances() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.ListDeploymentInstancesInput{
		DeploymentId: aws.String("DeploymentId"), // Required
		InstanceStatusFilter: []*string{
			aws.String("InstanceStatus"), // Required
			// More values...
		},
		InstanceTypeFilter: []*string{
			aws.String("InstanceType"), // Required
			// More values...
		},
		NextToken: aws.String("NextToken"),
	}
	resp, err := svc.ListDeploymentInstances(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_ListDeployments() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.ListDeploymentsInput{
		ApplicationName: aws.String("ApplicationName"),
		CreateTimeRange: &codedeploy.TimeRange{
			End:   aws.Time(time.Now()),
			Start: aws.Time(time.Now()),
		},
		DeploymentGroupName: aws.String("DeploymentGroupName"),
		IncludeOnlyStatuses: []*string{
			aws.String("DeploymentStatus"), // Required
			// More values...
		},
		NextToken: aws.String("NextToken"),
	}
	resp, err := svc.ListDeployments(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_ListOnPremisesInstances() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.ListOnPremisesInstancesInput{
		NextToken:          aws.String("NextToken"),
		RegistrationStatus: aws.String("RegistrationStatus"),
		TagFilters: []*codedeploy.TagFilter{
			{ // Required
				Key:   aws.String("Key"),
				Type:  aws.String("TagFilterType"),
				Value: aws.String("Value"),
			},
			// More values...
		},
	}
	resp, err := svc.ListOnPremisesInstances(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_RegisterApplicationRevision() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.RegisterApplicationRevisionInput{
		ApplicationName: aws.String("ApplicationName"), // Required
		Revision: &codedeploy.RevisionLocation{ // Required
			GitHubLocation: &codedeploy.GitHubLocation{
				CommitId:   aws.String("CommitId"),
				Repository: aws.String("Repository"),
			},
			RevisionType: aws.String("RevisionLocationType"),
			S3Location: &codedeploy.S3Location{
				Bucket:     aws.String("S3Bucket"),
				BundleType: aws.String("BundleType"),
				ETag:       aws.String("ETag"),
				Key:        aws.String("S3Key"),
				Version:    aws.String("VersionId"),
			},
		},
		Description: aws.String("Description"),
	}
	resp, err := svc.RegisterApplicationRevision(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_RegisterOnPremisesInstance() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.RegisterOnPremisesInstanceInput{
		InstanceName:  aws.String("InstanceName"), // Required
		IamSessionArn: aws.String("IamSessionArn"),
		IamUserArn:    aws.String("IamUserArn"),
	}
	resp, err := svc.RegisterOnPremisesInstance(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_RemoveTagsFromOnPremisesInstances() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.RemoveTagsFromOnPremisesInstancesInput{
		InstanceNames: []*string{ // Required
			aws.String("InstanceName"), // Required
			// More values...
		},
		Tags: []*codedeploy.Tag{ // Required
			{ // Required
				Key:   aws.String("Key"),
				Value: aws.String("Value"),
			},
			// More values...
		},
	}
	resp, err := svc.RemoveTagsFromOnPremisesInstances(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_SkipWaitTimeForInstanceTermination() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.SkipWaitTimeForInstanceTerminationInput{
		DeploymentId: aws.String("DeploymentId"),
	}
	resp, err := svc.SkipWaitTimeForInstanceTermination(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_StopDeployment() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.StopDeploymentInput{
		DeploymentId:        aws.String("DeploymentId"), // Required
		AutoRollbackEnabled: aws.Bool(true),
	}
	resp, err := svc.StopDeployment(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_UpdateApplication() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.UpdateApplicationInput{
		ApplicationName:    aws.String("ApplicationName"),
		NewApplicationName: aws.String("ApplicationName"),
	}
	resp, err := svc.UpdateApplication(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleCodeDeploy_UpdateDeploymentGroup() {
	sess := session.Must(session.NewSession())

	svc := codedeploy.New(sess)

	params := &codedeploy.UpdateDeploymentGroupInput{
		ApplicationName:            aws.String("ApplicationName"),     // Required
		CurrentDeploymentGroupName: aws.String("DeploymentGroupName"), // Required
		AlarmConfiguration: &codedeploy.AlarmConfiguration{
			Alarms: []*codedeploy.Alarm{
				{ // Required
					Name: aws.String("AlarmName"),
				},
				// More values...
			},
			Enabled:                aws.Bool(true),
			IgnorePollAlarmFailure: aws.Bool(true),
		},
		AutoRollbackConfiguration: &codedeploy.AutoRollbackConfiguration{
			Enabled: aws.Bool(true),
			Events: []*string{
				aws.String("AutoRollbackEvent"), // Required
				// More values...
			},
		},
		AutoScalingGroups: []*string{
			aws.String("AutoScalingGroupName"), // Required
			// More values...
		},
		BlueGreenDeploymentConfiguration: &codedeploy.BlueGreenDeploymentConfiguration{
			DeploymentReadyOption: &codedeploy.DeploymentReadyOption{
				ActionOnTimeout:   aws.String("DeploymentReadyAction"),
				WaitTimeInMinutes: aws.Int64(1),
			},
			GreenFleetProvisioningOption: &codedeploy.GreenFleetProvisioningOption{
				Action: aws.String("GreenFleetProvisioningAction"),
			},
			TerminateBlueInstancesOnDeploymentSuccess: &codedeploy.BlueInstanceTerminationOption{
				Action: aws.String("InstanceAction"),
				TerminationWaitTimeInMinutes: aws.Int64(1),
			},
		},
		DeploymentConfigName: aws.String("DeploymentConfigName"),
		DeploymentStyle: &codedeploy.DeploymentStyle{
			DeploymentOption: aws.String("DeploymentOption"),
			DeploymentType:   aws.String("DeploymentType"),
		},
		Ec2TagFilters: []*codedeploy.EC2TagFilter{
			{ // Required
				Key:   aws.String("Key"),
				Type:  aws.String("EC2TagFilterType"),
				Value: aws.String("Value"),
			},
			// More values...
		},
		LoadBalancerInfo: &codedeploy.LoadBalancerInfo{
			ElbInfoList: []*codedeploy.ELBInfo{
				{ // Required
					Name: aws.String("ELBName"),
				},
				// More values...
			},
		},
		NewDeploymentGroupName: aws.String("DeploymentGroupName"),
		OnPremisesInstanceTagFilters: []*codedeploy.TagFilter{
			{ // Required
				Key:   aws.String("Key"),
				Type:  aws.String("TagFilterType"),
				Value: aws.String("Value"),
			},
			// More values...
		},
		ServiceRoleArn: aws.String("Role"),
		TriggerConfigurations: []*codedeploy.TriggerConfig{
			{ // Required
				TriggerEvents: []*string{
					aws.String("TriggerEventType"), // Required
					// More values...
				},
				TriggerName:      aws.String("TriggerName"),
				TriggerTargetArn: aws.String("TriggerTargetArn"),
			},
			// More values...
		},
	}
	resp, err := svc.UpdateDeploymentGroup(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}
