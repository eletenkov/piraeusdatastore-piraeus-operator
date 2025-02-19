package merge

import (
	"sort"

	piraeusv1 "github.com/piraeusdatastore/piraeus-operator/v2/api/v1"
)

// SatelliteConfigurations merges all configurations that apply based on the given node labels
//
// Merging happens by:
// * Concatenating all patches in the matching configs
// * Merging all properties by name. A property defined in a "later" config overrides previous property definitions.
// * Merging all storage pools by name. A storage pool defined in a "later" config overrides previous property definitions.
func SatelliteConfigurations(nodeLabels map[string]string, configs ...piraeusv1.LinstorSatelliteConfiguration) *piraeusv1.LinstorSatelliteConfiguration {
	result := &piraeusv1.LinstorSatelliteConfiguration{}

	propsMap := make(map[string]*piraeusv1.LinstorNodeProperty)
	storPoolMap := make(map[string]*piraeusv1.LinstorStoragePool)

	for i := range configs {
		cfg := &configs[i]

		if !SubsetOf(cfg.Spec.NodeSelector, nodeLabels) {
			continue
		}

		for j := range cfg.Spec.Properties {
			propsMap[cfg.Spec.Properties[j].Name] = &cfg.Spec.Properties[j]
		}

		for j := range cfg.Spec.StoragePools {
			storPoolMap[cfg.Spec.StoragePools[j].Name] = &cfg.Spec.StoragePools[j]
		}

		result.Spec.Patches = append(result.Spec.Patches, cfg.Spec.Patches...)

		if cfg.Spec.InternalTLS != nil {
			result.Spec.InternalTLS = cfg.Spec.InternalTLS
		}
	}

	for _, v := range propsMap {
		result.Spec.Properties = append(result.Spec.Properties, *v)
	}

	sort.Slice(result.Spec.Properties, func(i, j int) bool {
		return result.Spec.Properties[i].Name < result.Spec.Properties[j].Name
	})

	for _, v := range storPoolMap {
		result.Spec.StoragePools = append(result.Spec.StoragePools, *v)
	}

	sort.Slice(result.Spec.StoragePools, func(i, j int) bool {
		return result.Spec.StoragePools[i].Name < result.Spec.StoragePools[j].Name
	})

	return result
}

// SubsetOf returns true if all key and values in sub also appear in super.
func SubsetOf(sub, super map[string]string) bool {
	for k, subv := range sub {
		superv, ok := super[k]
		if !ok || superv != subv {
			return false
		}
	}

	return true
}
