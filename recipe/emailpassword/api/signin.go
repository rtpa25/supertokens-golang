package api

import (
	"encoding/json"
	"io/ioutil"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func SignInAPI(apiImplementation epmodels.APIInterface, options epmodels.APIOptions) error {
	if apiImplementation.SignInPOST == nil {
		options.OtherHandler(options.Res, options.Req)
		return nil
	}

	body, err := ioutil.ReadAll(options.Req.Body)
	if err != nil {
		return err
	}
	var formFieldsRaw map[string]interface{}
	err = json.Unmarshal(body, &formFieldsRaw)
	if err != nil {
		return err
	}

	formFields, err := validateFormFieldsOrThrowError(options.Config.SignInFeature.FormFields, formFieldsRaw["formFields"].([]interface{}))
	if err != nil {
		return err
	}

	result, err := apiImplementation.SignInPOST(formFields, options)
	if err != nil {
		return err
	}
	if result.WrongCredentialsError != nil {
		return supertokens.Send200Response(options.Res, map[string]interface{}{
			"status": "WRONG_CREDENTIALS_ERROR",
		})
	} else {
		return supertokens.Send200Response(options.Res, map[string]interface{}{
			"status": "OK",
			"user":   result.OK.User,
		})
	}
}
