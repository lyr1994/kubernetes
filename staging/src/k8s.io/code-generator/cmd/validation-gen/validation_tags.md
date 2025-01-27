# Kubernetes Validation Tags Documentation

This document lists the supported validation tags and their related information.

## Tags Overview

| Tag | Description | Contexts |
|-----|-------------|----------|
| [`k8s:eachKey`](#k8seachkey) | Declares a validation for map keys. | Type definition, Field definition, map key, map/slice value |
| [`k8s:eachVal`](#k8seachval) | Declares a validation for map and slice values. | Type definition, Field definition, map key, map/slice value |
| `k8s:enum` | Indicates that a string type is an enum. All const values of this type are considered values in the enum. | Type definition |
| `k8s:forbidden` | Indicates that a field is forbidden to be specified. | Type definition, Field definition, map key, map/slice value |
| [`k8s:format`](#k8sformat) | Indicates that a string field has a particular format. | Type definition, Field definition, map key, map/slice value |
| [`k8s:ifOptionDisabled(<option-name>)`](#k8sifoptiondisabled(<option-name>)) | Declares a validation that only applies when an option is disabled. | Type definition, Field definition, map key, map/slice value |
| [`k8s:ifOptionEnabled(<option-name>)`](#k8sifoptionenabled(<option-name>)) | Declares a validation that only applies when an option is enabled. | Type definition, Field definition, map key, map/slice value |
| [`k8s:listMapKey`](#k8slistmapkey) | Declares a named field of a list's value type as part of the list-map key. | Type definition, Field definition, map key, map/slice value |
| [`k8s:maxItems`](#k8smaxitems) | Indidates that a slice field has a limit on its size. | Type definition, Field definition, map key, map/slice value |
| [`k8s:maxLength`](#k8smaxlength) | Indicates that a string field has a limit on its length. | Type definition, Field definition, map key, map/slice value |
| `k8s:optional` | Indicates that a field is optional. | Type definition, Field definition, map key, map/slice value |
| `k8s:required` | Indicates that a field is required to be specified. | Type definition, Field definition, map key, map/slice value |
| [`k8s:subfield(<field-name>)`](#k8ssubfield(<field-name>)) | Declares a validation for a specified subfield of the struct. The subfield must be a direct field of the struct, or of an embedded struct | Field definition, map key, map/slice value |
| [`k8s:unionDiscriminator`](#k8suniondiscriminator) | Indicates that this field is the discriminator for a union. | Field definition, map key, map/slice value |
| [`k8s:unionMember`](#k8sunionmember) | Indicates that this field is a member of a union. | Field definition, map key, map/slice value |
| [`k8s:validateError`](#k8svalidateerror) | Always fails code generation (useful for testing). | Type definition, Field definition, map key, map/slice value |
| [`k8s:validateFalse`](#k8svalidatefalse) | Always fails validation (useful for testing). | Type definition, Field definition, map key, map/slice value |
| [`k8s:validateTrue`](#k8svalidatetrue) | Always passes validation (useful for testing). | Type definition, Field definition, map key, map/slice value |

## Tag Details

### k8s:eachKey

| Description | Docs | Schema |
|-------------|------|---------|
| **\<validation-tag\>** | This tag will be evaluated for each key of a map. | None |

### k8s:eachVal

| Description | Docs | Schema |
|-------------|------|---------|
| **\<validation-tag\>** | This tag will be evaluated for each value of a map or slice. | None |

### k8s:format

| Description | Docs | Schema |
|-------------|------|---------|
| **ip-sloppy** | This field holds an IPv4 or IPv6 address value. IPv4 octets may have leading zeros. | None |
| **dns-label** | This field holds a DNS label value. | None |

### k8s:ifOptionDisabled(<option-name>)

| Description | Docs | Schema |
|-------------|------|---------|
| **\<validation-tag\>** | This validation tag will be evaluated only if the validation option is disabled. | None |

### k8s:ifOptionEnabled(<option-name>)

| Description | Docs | Schema |
|-------------|------|---------|
| **\<validation-tag\>** | This validation tag will be evaluated only if the validation option is enabled. | None |

### k8s:listMapKey

| Description | Docs | Schema |
|-------------|------|---------|
| **\<field-name\>** | This values names a field of a list's value type. | None |

### k8s:maxItems

| Description | Docs | Schema |
|-------------|------|---------|
| **\<non-negative integer\>** | This field must be no more than X items long. | None |

### k8s:maxLength

| Description | Docs | Schema |
|-------------|------|---------|
| **\<non-negative integer\>** | This field must be no more than X characters long. | None |

### k8s:subfield(<field-name>)

| Description | Docs | Schema |
|-------------|------|---------|
| **\<validation-tag\>** | This tag will be evaluated for the subfield of the struct. | None |

### k8s:unionDiscriminator

| Description | Docs | Schema |
|-------------|------|---------|
| **\<json-object\>** |  | - `union`: `<string>` (the name of the union, if more than one exists) |

### k8s:unionMember

| Description | Docs | Schema |
|-------------|------|---------|
| **\<json-object\>** |  | - `union`: `<string>` (the name of the union, if more than one exists)<br>- `memberName`: `<string>` (the discriminator value for this member) |

### k8s:validateError

| Description | Docs | Schema |
|-------------|------|---------|
| **\<string\>** | This string will be included in the error message. | None |

### k8s:validateFalse

| Description | Docs | Schema |
|-------------|------|---------|
| **\<none\>** |  | None |
| **\<quoted-string\>** | The generated code will include this string. | None |
| **\<json-object\>** |  | - `flags`: `<list-of-flag-string>` (values: ShortCircuit, NonError)<br>- `msg`: `<string>` (The generated code will include this string.) |

### k8s:validateTrue

| Description | Docs | Schema |
|-------------|------|---------|
| **\<none\>** |  | None |
| **\<quoted-string\>** | The generated code will include this string. | None |
| **\<json-object\>** |  | - `flags`: `<list-of-flag-string>` (values: ShortCircuit, NonError)<br>- `msg`: `<string>` (The generated code will include this string.)<br>- `typeArg`: `<string>` (The type arg in generated code (must be the value-type, not pointer).) |

