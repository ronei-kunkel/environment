package environment

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

func Load[T any](sources ...string) T {
	var reflectErrs []error

	dotEnvErr := godotenv.Load(sources...)

	var instance T

	v := reflect.ValueOf(&instance).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		envVar := field.Tag.Get("env")
		if envVar == "" {
			envVar = field.Name
		}
		value := os.Getenv(envVar)
		if value == "" {
			reflectErrs = append(reflectErrs, fmt.Errorf("has no `%s` environment variable defined to populate into `%s` instance field", envVar, field.Name))
		}
		v.Field(i).SetString(value)
	}

	if len(reflectErrs) == 0 {
		return instance
	}

	if dotEnvErr != nil {
		reflectErrs = append(reflectErrs, fmt.Errorf("%s", dotEnvErr.Error()))
	}

	log.Println("Errors loading environment variables:")
	for _, e := range reflectErrs {
		log.Println(" -", e)
	}
	log.Fatalln("Aborting due to missing env vars")

	return instance
}
