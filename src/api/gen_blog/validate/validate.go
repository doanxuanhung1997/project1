package validate

import (
	"github.com/go-playground/validator/v10"
	"houze_ops_backend/api/gen_blog/model"
)

var (
	validate *validator.Validate
)

type ApiError struct {
	Param   string
	Message string
}

func InputCreateBlogValidate(input model.InputCreateBlog) []ApiError {
	validate = validator.New()
	err := validate.Struct(input)
	if err != nil {
		out := make([]ApiError, len(err.(validator.ValidationErrors)))
		for i, fe := range err.(validator.ValidationErrors) {
			out[i] = ApiError{fe.Field(), msgForTag(fe)}
		}
		//for _, err := range err.(validator.ValidationErrors) {
		//
		//	fmt.Println(err.Namespace())
		//	fmt.Println(err.Field())
		//	fmt.Println(err.StructNamespace())
		//	fmt.Println(err.StructField())
		//	fmt.Println(err.Tag())
		//	fmt.Println(err.ActualTag())
		//	fmt.Println(err.Kind())
		//	fmt.Println(err.Type())
		//	fmt.Println(err.Value())
		//	fmt.Println(err.Param())
		//	fmt.Println()
		//}

		return out
	}
	return nil
}

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	}
	return fe.Error() // default error
}
