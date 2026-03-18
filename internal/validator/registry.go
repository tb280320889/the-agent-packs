package validator

func registry() map[string]ValidatorFunc {
	return map[string]ValidatorFunc{
		"validator-core-output":         validateCoreOutput,
		"validator-contract-delivery":   validateContractDelivery,
		"validator-domain-wxt-manifest": validateDomainWXTManifest,
	}
}
