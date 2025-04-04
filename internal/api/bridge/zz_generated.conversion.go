//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by conversion-gen. DO NOT EDIT.

package bridge

import (
	unsafe "unsafe"

	v1alpha1 "github.com/aws/eks-hybrid/api/v1alpha1"
	api "github.com/aws/eks-hybrid/internal/api"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*v1alpha1.ClusterDetails)(nil), (*api.ClusterDetails)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ClusterDetails_To_api_ClusterDetails(a.(*v1alpha1.ClusterDetails), b.(*api.ClusterDetails), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*api.ClusterDetails)(nil), (*v1alpha1.ClusterDetails)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_api_ClusterDetails_To_v1alpha1_ClusterDetails(a.(*api.ClusterDetails), b.(*v1alpha1.ClusterDetails), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.ContainerdOptions)(nil), (*api.ContainerdOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ContainerdOptions_To_api_ContainerdOptions(a.(*v1alpha1.ContainerdOptions), b.(*api.ContainerdOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*api.ContainerdOptions)(nil), (*v1alpha1.ContainerdOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_api_ContainerdOptions_To_v1alpha1_ContainerdOptions(a.(*api.ContainerdOptions), b.(*v1alpha1.ContainerdOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.HybridOptions)(nil), (*api.HybridOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_HybridOptions_To_api_HybridOptions(a.(*v1alpha1.HybridOptions), b.(*api.HybridOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*api.HybridOptions)(nil), (*v1alpha1.HybridOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_api_HybridOptions_To_v1alpha1_HybridOptions(a.(*api.HybridOptions), b.(*v1alpha1.HybridOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.IAMRolesAnywhere)(nil), (*api.IAMRolesAnywhere)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_IAMRolesAnywhere_To_api_IAMRolesAnywhere(a.(*v1alpha1.IAMRolesAnywhere), b.(*api.IAMRolesAnywhere), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*api.IAMRolesAnywhere)(nil), (*v1alpha1.IAMRolesAnywhere)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_api_IAMRolesAnywhere_To_v1alpha1_IAMRolesAnywhere(a.(*api.IAMRolesAnywhere), b.(*v1alpha1.IAMRolesAnywhere), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.InstanceOptions)(nil), (*api.InstanceOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_InstanceOptions_To_api_InstanceOptions(a.(*v1alpha1.InstanceOptions), b.(*api.InstanceOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*api.InstanceOptions)(nil), (*v1alpha1.InstanceOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_api_InstanceOptions_To_v1alpha1_InstanceOptions(a.(*api.InstanceOptions), b.(*v1alpha1.InstanceOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.KubeletOptions)(nil), (*api.KubeletOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_KubeletOptions_To_api_KubeletOptions(a.(*v1alpha1.KubeletOptions), b.(*api.KubeletOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*api.KubeletOptions)(nil), (*v1alpha1.KubeletOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_api_KubeletOptions_To_v1alpha1_KubeletOptions(a.(*api.KubeletOptions), b.(*v1alpha1.KubeletOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.LocalStorageOptions)(nil), (*api.LocalStorageOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_LocalStorageOptions_To_api_LocalStorageOptions(a.(*v1alpha1.LocalStorageOptions), b.(*api.LocalStorageOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*api.LocalStorageOptions)(nil), (*v1alpha1.LocalStorageOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_api_LocalStorageOptions_To_v1alpha1_LocalStorageOptions(a.(*api.LocalStorageOptions), b.(*v1alpha1.LocalStorageOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.NodeConfig)(nil), (*api.NodeConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_NodeConfig_To_api_NodeConfig(a.(*v1alpha1.NodeConfig), b.(*api.NodeConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*api.NodeConfig)(nil), (*v1alpha1.NodeConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_api_NodeConfig_To_v1alpha1_NodeConfig(a.(*api.NodeConfig), b.(*v1alpha1.NodeConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.NodeConfigList)(nil), (*api.NodeConfigList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_NodeConfigList_To_api_NodeConfigList(a.(*v1alpha1.NodeConfigList), b.(*api.NodeConfigList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*api.NodeConfigList)(nil), (*v1alpha1.NodeConfigList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_api_NodeConfigList_To_v1alpha1_NodeConfigList(a.(*api.NodeConfigList), b.(*v1alpha1.NodeConfigList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.NodeConfigSpec)(nil), (*api.NodeConfigSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_NodeConfigSpec_To_api_NodeConfigSpec(a.(*v1alpha1.NodeConfigSpec), b.(*api.NodeConfigSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*api.NodeConfigSpec)(nil), (*v1alpha1.NodeConfigSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_api_NodeConfigSpec_To_v1alpha1_NodeConfigSpec(a.(*api.NodeConfigSpec), b.(*v1alpha1.NodeConfigSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.SSM)(nil), (*api.SSM)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_SSM_To_api_SSM(a.(*v1alpha1.SSM), b.(*api.SSM), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*api.SSM)(nil), (*v1alpha1.SSM)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_api_SSM_To_v1alpha1_SSM(a.(*api.SSM), b.(*v1alpha1.SSM), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1alpha1_ClusterDetails_To_api_ClusterDetails(in *v1alpha1.ClusterDetails, out *api.ClusterDetails, s conversion.Scope) error {
	out.Name = in.Name
	out.Region = in.Region
	out.APIServerEndpoint = in.APIServerEndpoint
	out.CertificateAuthority = *(*[]byte)(unsafe.Pointer(&in.CertificateAuthority))
	out.CIDR = in.CIDR
	out.EnableOutpost = (*bool)(unsafe.Pointer(in.EnableOutpost))
	out.ID = in.ID
	return nil
}

// Convert_v1alpha1_ClusterDetails_To_api_ClusterDetails is an autogenerated conversion function.
func Convert_v1alpha1_ClusterDetails_To_api_ClusterDetails(in *v1alpha1.ClusterDetails, out *api.ClusterDetails, s conversion.Scope) error {
	return autoConvert_v1alpha1_ClusterDetails_To_api_ClusterDetails(in, out, s)
}

func autoConvert_api_ClusterDetails_To_v1alpha1_ClusterDetails(in *api.ClusterDetails, out *v1alpha1.ClusterDetails, s conversion.Scope) error {
	out.Name = in.Name
	out.Region = in.Region
	out.APIServerEndpoint = in.APIServerEndpoint
	out.CertificateAuthority = *(*[]byte)(unsafe.Pointer(&in.CertificateAuthority))
	out.CIDR = in.CIDR
	out.EnableOutpost = (*bool)(unsafe.Pointer(in.EnableOutpost))
	out.ID = in.ID
	return nil
}

// Convert_api_ClusterDetails_To_v1alpha1_ClusterDetails is an autogenerated conversion function.
func Convert_api_ClusterDetails_To_v1alpha1_ClusterDetails(in *api.ClusterDetails, out *v1alpha1.ClusterDetails, s conversion.Scope) error {
	return autoConvert_api_ClusterDetails_To_v1alpha1_ClusterDetails(in, out, s)
}

func autoConvert_v1alpha1_ContainerdOptions_To_api_ContainerdOptions(in *v1alpha1.ContainerdOptions, out *api.ContainerdOptions, s conversion.Scope) error {
	out.Config = in.Config
	return nil
}

// Convert_v1alpha1_ContainerdOptions_To_api_ContainerdOptions is an autogenerated conversion function.
func Convert_v1alpha1_ContainerdOptions_To_api_ContainerdOptions(in *v1alpha1.ContainerdOptions, out *api.ContainerdOptions, s conversion.Scope) error {
	return autoConvert_v1alpha1_ContainerdOptions_To_api_ContainerdOptions(in, out, s)
}

func autoConvert_api_ContainerdOptions_To_v1alpha1_ContainerdOptions(in *api.ContainerdOptions, out *v1alpha1.ContainerdOptions, s conversion.Scope) error {
	out.Config = in.Config
	return nil
}

// Convert_api_ContainerdOptions_To_v1alpha1_ContainerdOptions is an autogenerated conversion function.
func Convert_api_ContainerdOptions_To_v1alpha1_ContainerdOptions(in *api.ContainerdOptions, out *v1alpha1.ContainerdOptions, s conversion.Scope) error {
	return autoConvert_api_ContainerdOptions_To_v1alpha1_ContainerdOptions(in, out, s)
}

func autoConvert_v1alpha1_HybridOptions_To_api_HybridOptions(in *v1alpha1.HybridOptions, out *api.HybridOptions, s conversion.Scope) error {
	out.EnableCredentialsFile = in.EnableCredentialsFile
	out.IAMRolesAnywhere = (*api.IAMRolesAnywhere)(unsafe.Pointer(in.IAMRolesAnywhere))
	out.SSM = (*api.SSM)(unsafe.Pointer(in.SSM))
	return nil
}

// Convert_v1alpha1_HybridOptions_To_api_HybridOptions is an autogenerated conversion function.
func Convert_v1alpha1_HybridOptions_To_api_HybridOptions(in *v1alpha1.HybridOptions, out *api.HybridOptions, s conversion.Scope) error {
	return autoConvert_v1alpha1_HybridOptions_To_api_HybridOptions(in, out, s)
}

func autoConvert_api_HybridOptions_To_v1alpha1_HybridOptions(in *api.HybridOptions, out *v1alpha1.HybridOptions, s conversion.Scope) error {
	out.EnableCredentialsFile = in.EnableCredentialsFile
	out.IAMRolesAnywhere = (*v1alpha1.IAMRolesAnywhere)(unsafe.Pointer(in.IAMRolesAnywhere))
	out.SSM = (*v1alpha1.SSM)(unsafe.Pointer(in.SSM))
	return nil
}

// Convert_api_HybridOptions_To_v1alpha1_HybridOptions is an autogenerated conversion function.
func Convert_api_HybridOptions_To_v1alpha1_HybridOptions(in *api.HybridOptions, out *v1alpha1.HybridOptions, s conversion.Scope) error {
	return autoConvert_api_HybridOptions_To_v1alpha1_HybridOptions(in, out, s)
}

func autoConvert_v1alpha1_IAMRolesAnywhere_To_api_IAMRolesAnywhere(in *v1alpha1.IAMRolesAnywhere, out *api.IAMRolesAnywhere, s conversion.Scope) error {
	out.NodeName = in.NodeName
	out.TrustAnchorARN = in.TrustAnchorARN
	out.ProfileARN = in.ProfileARN
	out.RoleARN = in.RoleARN
	out.AwsConfigPath = in.AwsConfigPath
	out.CertificatePath = in.CertificatePath
	out.PrivateKeyPath = in.PrivateKeyPath
	return nil
}

// Convert_v1alpha1_IAMRolesAnywhere_To_api_IAMRolesAnywhere is an autogenerated conversion function.
func Convert_v1alpha1_IAMRolesAnywhere_To_api_IAMRolesAnywhere(in *v1alpha1.IAMRolesAnywhere, out *api.IAMRolesAnywhere, s conversion.Scope) error {
	return autoConvert_v1alpha1_IAMRolesAnywhere_To_api_IAMRolesAnywhere(in, out, s)
}

func autoConvert_api_IAMRolesAnywhere_To_v1alpha1_IAMRolesAnywhere(in *api.IAMRolesAnywhere, out *v1alpha1.IAMRolesAnywhere, s conversion.Scope) error {
	out.NodeName = in.NodeName
	out.TrustAnchorARN = in.TrustAnchorARN
	out.ProfileARN = in.ProfileARN
	out.RoleARN = in.RoleARN
	out.AwsConfigPath = in.AwsConfigPath
	out.CertificatePath = in.CertificatePath
	out.PrivateKeyPath = in.PrivateKeyPath
	return nil
}

// Convert_api_IAMRolesAnywhere_To_v1alpha1_IAMRolesAnywhere is an autogenerated conversion function.
func Convert_api_IAMRolesAnywhere_To_v1alpha1_IAMRolesAnywhere(in *api.IAMRolesAnywhere, out *v1alpha1.IAMRolesAnywhere, s conversion.Scope) error {
	return autoConvert_api_IAMRolesAnywhere_To_v1alpha1_IAMRolesAnywhere(in, out, s)
}

func autoConvert_v1alpha1_InstanceOptions_To_api_InstanceOptions(in *v1alpha1.InstanceOptions, out *api.InstanceOptions, s conversion.Scope) error {
	if err := Convert_v1alpha1_LocalStorageOptions_To_api_LocalStorageOptions(&in.LocalStorage, &out.LocalStorage, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_InstanceOptions_To_api_InstanceOptions is an autogenerated conversion function.
func Convert_v1alpha1_InstanceOptions_To_api_InstanceOptions(in *v1alpha1.InstanceOptions, out *api.InstanceOptions, s conversion.Scope) error {
	return autoConvert_v1alpha1_InstanceOptions_To_api_InstanceOptions(in, out, s)
}

func autoConvert_api_InstanceOptions_To_v1alpha1_InstanceOptions(in *api.InstanceOptions, out *v1alpha1.InstanceOptions, s conversion.Scope) error {
	if err := Convert_api_LocalStorageOptions_To_v1alpha1_LocalStorageOptions(&in.LocalStorage, &out.LocalStorage, s); err != nil {
		return err
	}
	return nil
}

// Convert_api_InstanceOptions_To_v1alpha1_InstanceOptions is an autogenerated conversion function.
func Convert_api_InstanceOptions_To_v1alpha1_InstanceOptions(in *api.InstanceOptions, out *v1alpha1.InstanceOptions, s conversion.Scope) error {
	return autoConvert_api_InstanceOptions_To_v1alpha1_InstanceOptions(in, out, s)
}

func autoConvert_v1alpha1_KubeletOptions_To_api_KubeletOptions(in *v1alpha1.KubeletOptions, out *api.KubeletOptions, s conversion.Scope) error {
	out.Config = *(*api.InlineDocument)(unsafe.Pointer(&in.Config))
	out.Flags = *(*[]string)(unsafe.Pointer(&in.Flags))
	return nil
}

// Convert_v1alpha1_KubeletOptions_To_api_KubeletOptions is an autogenerated conversion function.
func Convert_v1alpha1_KubeletOptions_To_api_KubeletOptions(in *v1alpha1.KubeletOptions, out *api.KubeletOptions, s conversion.Scope) error {
	return autoConvert_v1alpha1_KubeletOptions_To_api_KubeletOptions(in, out, s)
}

func autoConvert_api_KubeletOptions_To_v1alpha1_KubeletOptions(in *api.KubeletOptions, out *v1alpha1.KubeletOptions, s conversion.Scope) error {
	out.Config = *(*map[string]runtime.RawExtension)(unsafe.Pointer(&in.Config))
	out.Flags = *(*[]string)(unsafe.Pointer(&in.Flags))
	return nil
}

// Convert_api_KubeletOptions_To_v1alpha1_KubeletOptions is an autogenerated conversion function.
func Convert_api_KubeletOptions_To_v1alpha1_KubeletOptions(in *api.KubeletOptions, out *v1alpha1.KubeletOptions, s conversion.Scope) error {
	return autoConvert_api_KubeletOptions_To_v1alpha1_KubeletOptions(in, out, s)
}

func autoConvert_v1alpha1_LocalStorageOptions_To_api_LocalStorageOptions(in *v1alpha1.LocalStorageOptions, out *api.LocalStorageOptions, s conversion.Scope) error {
	out.Strategy = api.LocalStorageStrategy(in.Strategy)
	return nil
}

// Convert_v1alpha1_LocalStorageOptions_To_api_LocalStorageOptions is an autogenerated conversion function.
func Convert_v1alpha1_LocalStorageOptions_To_api_LocalStorageOptions(in *v1alpha1.LocalStorageOptions, out *api.LocalStorageOptions, s conversion.Scope) error {
	return autoConvert_v1alpha1_LocalStorageOptions_To_api_LocalStorageOptions(in, out, s)
}

func autoConvert_api_LocalStorageOptions_To_v1alpha1_LocalStorageOptions(in *api.LocalStorageOptions, out *v1alpha1.LocalStorageOptions, s conversion.Scope) error {
	out.Strategy = v1alpha1.LocalStorageStrategy(in.Strategy)
	return nil
}

// Convert_api_LocalStorageOptions_To_v1alpha1_LocalStorageOptions is an autogenerated conversion function.
func Convert_api_LocalStorageOptions_To_v1alpha1_LocalStorageOptions(in *api.LocalStorageOptions, out *v1alpha1.LocalStorageOptions, s conversion.Scope) error {
	return autoConvert_api_LocalStorageOptions_To_v1alpha1_LocalStorageOptions(in, out, s)
}

func autoConvert_v1alpha1_NodeConfig_To_api_NodeConfig(in *v1alpha1.NodeConfig, out *api.NodeConfig, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1alpha1_NodeConfigSpec_To_api_NodeConfigSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_NodeConfig_To_api_NodeConfig is an autogenerated conversion function.
func Convert_v1alpha1_NodeConfig_To_api_NodeConfig(in *v1alpha1.NodeConfig, out *api.NodeConfig, s conversion.Scope) error {
	return autoConvert_v1alpha1_NodeConfig_To_api_NodeConfig(in, out, s)
}

func autoConvert_api_NodeConfig_To_v1alpha1_NodeConfig(in *api.NodeConfig, out *v1alpha1.NodeConfig, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_api_NodeConfigSpec_To_v1alpha1_NodeConfigSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	// INFO: in.Status opted out of conversion generation
	return nil
}

// Convert_api_NodeConfig_To_v1alpha1_NodeConfig is an autogenerated conversion function.
func Convert_api_NodeConfig_To_v1alpha1_NodeConfig(in *api.NodeConfig, out *v1alpha1.NodeConfig, s conversion.Scope) error {
	return autoConvert_api_NodeConfig_To_v1alpha1_NodeConfig(in, out, s)
}

func autoConvert_v1alpha1_NodeConfigList_To_api_NodeConfigList(in *v1alpha1.NodeConfigList, out *api.NodeConfigList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]api.NodeConfig, len(*in))
		for i := range *in {
			if err := Convert_v1alpha1_NodeConfig_To_api_NodeConfig(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

// Convert_v1alpha1_NodeConfigList_To_api_NodeConfigList is an autogenerated conversion function.
func Convert_v1alpha1_NodeConfigList_To_api_NodeConfigList(in *v1alpha1.NodeConfigList, out *api.NodeConfigList, s conversion.Scope) error {
	return autoConvert_v1alpha1_NodeConfigList_To_api_NodeConfigList(in, out, s)
}

func autoConvert_api_NodeConfigList_To_v1alpha1_NodeConfigList(in *api.NodeConfigList, out *v1alpha1.NodeConfigList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]v1alpha1.NodeConfig, len(*in))
		for i := range *in {
			if err := Convert_api_NodeConfig_To_v1alpha1_NodeConfig(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

// Convert_api_NodeConfigList_To_v1alpha1_NodeConfigList is an autogenerated conversion function.
func Convert_api_NodeConfigList_To_v1alpha1_NodeConfigList(in *api.NodeConfigList, out *v1alpha1.NodeConfigList, s conversion.Scope) error {
	return autoConvert_api_NodeConfigList_To_v1alpha1_NodeConfigList(in, out, s)
}

func autoConvert_v1alpha1_NodeConfigSpec_To_api_NodeConfigSpec(in *v1alpha1.NodeConfigSpec, out *api.NodeConfigSpec, s conversion.Scope) error {
	if err := Convert_v1alpha1_ClusterDetails_To_api_ClusterDetails(&in.Cluster, &out.Cluster, s); err != nil {
		return err
	}
	if err := Convert_v1alpha1_ContainerdOptions_To_api_ContainerdOptions(&in.Containerd, &out.Containerd, s); err != nil {
		return err
	}
	if err := Convert_v1alpha1_InstanceOptions_To_api_InstanceOptions(&in.Instance, &out.Instance, s); err != nil {
		return err
	}
	if err := Convert_v1alpha1_KubeletOptions_To_api_KubeletOptions(&in.Kubelet, &out.Kubelet, s); err != nil {
		return err
	}
	out.Hybrid = (*api.HybridOptions)(unsafe.Pointer(in.Hybrid))
	return nil
}

// Convert_v1alpha1_NodeConfigSpec_To_api_NodeConfigSpec is an autogenerated conversion function.
func Convert_v1alpha1_NodeConfigSpec_To_api_NodeConfigSpec(in *v1alpha1.NodeConfigSpec, out *api.NodeConfigSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_NodeConfigSpec_To_api_NodeConfigSpec(in, out, s)
}

func autoConvert_api_NodeConfigSpec_To_v1alpha1_NodeConfigSpec(in *api.NodeConfigSpec, out *v1alpha1.NodeConfigSpec, s conversion.Scope) error {
	if err := Convert_api_ClusterDetails_To_v1alpha1_ClusterDetails(&in.Cluster, &out.Cluster, s); err != nil {
		return err
	}
	if err := Convert_api_ContainerdOptions_To_v1alpha1_ContainerdOptions(&in.Containerd, &out.Containerd, s); err != nil {
		return err
	}
	if err := Convert_api_InstanceOptions_To_v1alpha1_InstanceOptions(&in.Instance, &out.Instance, s); err != nil {
		return err
	}
	if err := Convert_api_KubeletOptions_To_v1alpha1_KubeletOptions(&in.Kubelet, &out.Kubelet, s); err != nil {
		return err
	}
	out.Hybrid = (*v1alpha1.HybridOptions)(unsafe.Pointer(in.Hybrid))
	return nil
}

// Convert_api_NodeConfigSpec_To_v1alpha1_NodeConfigSpec is an autogenerated conversion function.
func Convert_api_NodeConfigSpec_To_v1alpha1_NodeConfigSpec(in *api.NodeConfigSpec, out *v1alpha1.NodeConfigSpec, s conversion.Scope) error {
	return autoConvert_api_NodeConfigSpec_To_v1alpha1_NodeConfigSpec(in, out, s)
}

func autoConvert_v1alpha1_SSM_To_api_SSM(in *v1alpha1.SSM, out *api.SSM, s conversion.Scope) error {
	out.ActivationCode = in.ActivationCode
	out.ActivationID = in.ActivationID
	return nil
}

// Convert_v1alpha1_SSM_To_api_SSM is an autogenerated conversion function.
func Convert_v1alpha1_SSM_To_api_SSM(in *v1alpha1.SSM, out *api.SSM, s conversion.Scope) error {
	return autoConvert_v1alpha1_SSM_To_api_SSM(in, out, s)
}

func autoConvert_api_SSM_To_v1alpha1_SSM(in *api.SSM, out *v1alpha1.SSM, s conversion.Scope) error {
	out.ActivationCode = in.ActivationCode
	out.ActivationID = in.ActivationID
	return nil
}

// Convert_api_SSM_To_v1alpha1_SSM is an autogenerated conversion function.
func Convert_api_SSM_To_v1alpha1_SSM(in *api.SSM, out *v1alpha1.SSM, s conversion.Scope) error {
	return autoConvert_api_SSM_To_v1alpha1_SSM(in, out, s)
}
