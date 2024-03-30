# REST API Conventions

Many of the conventions used in the example microservice's REST APIs follow guidance from [restfulapi.net](https://restfulapi.net/resource-naming/)
but deviate slightly by allowing RPC-like actions off of entities and collections presented in the API.

## Casing and versioning

APIs are always lowercase and follow kebab casing. The following URLs are OK:

* `/api/v1/countries/usa/states/maryland/population-by-hair-color`
* `/api/v1/countries/usa/states?sort-by=name&sort-direction=descending`
* `/api/v1/countries/1/states/5`

These URLs are not:

* `/api/v1/countries/usa/states/maryland/populationByHairColor` - Don't use camel case
* `/api/v1/countries/usa/states?sort_by=name&sort_direction=descending` - Don't use snake case

This naming rule does not apply to values passed to query parameters, so the following would be okay:

* `/api/v1/countries?name-like=United%20States%20Of%20America`

The version of a specific API endpoint is included in the URL after the "/api" prefix. It should be incremented if 
breaking changes are made to the endpoint, and the old endpoint should continue to be served. The old endpoint should 
follow a deprecation cycle across a 2-month period if it is to be decommissioned in favor of the new endpoint. That is,
the old endpoint should be marked deprecated in the swagger docs, migration instructions should be added, and these 
should be announced in an upcoming release announcement with a second reminder in the release announcement preceding 
decommissioning.

Breaking changes are defined as:

* Endpoint naming changes
* Removing fields from a response body
* Adding required information in a request
* Renaming fields in request or response bodies

These changes are not considered breaking:

* Adding new fields to a response body (consumers just won't deserialize those fields)
* Removing required information in a request (though this change should be announced)

## Collections

Collections of resources should always be **nouns with plural names**.

| Do This              | Don't do This      |
|----------------------|--------------------|
| `/api/v1/countries/1` | `/api/v1/country/1` |
| `/api/v1/countries/1` | `/api/v1/get-country/1` |


Exposing collections of resources like this clearly states to API consumers that the resource has more than one entry
and could be filtered or may have sub-resources that can be extracted.

The HTTP verb `GET` should be used to list the entries in the collection, and query parameters may be used to filter,
sort, or paginate the contents.

Example: `GET /api/v1/countries?sort-by=flagColor&page=2`

Retrieving the contents of a collection should return a `200 OK`. If the collection is empty, a `200 OK` should be sent
along with an empty array in the response body. `404 Not Found` should only be triggered on a collection if it's a collection
nested inside another resource, and the specific resource wasn't found.

For example, a request to `GET /api/v1/users/1/friends` should return a `404 NOT FOUND` if user 1 doesn't exist. The response
body should then indicate that specifically the user was the missing piece.

### Adding to the collection

Entries can be added to the collection by using the HTTP verb `POST`. `PUT` may be used to perform an upsert 
(update if exists, otherwise insert) rather than an insertion into the collection. Requests made with these verbs
contain the contents of the new entry in the request body.

Fields in request bodies should be camel cased (via the go "json" struct tag, see [the microservice architecture docs](Microservice Architecture.md#dtos-validation-and-responding)
for more info) and acronyms should only capitalize the first letter.

Examples:

```http request
POST /api/v1/countries
Content-Type: application/json

{
  "name": "Genovia",
  "flagUrl": "https://en.wikipedia.org/wiki/File:Flag_of_Genovia.svg"
}
```

```http request
PUT /api/v1/countries/usa/states/maryland/state-icons
Content-Type: application/json

{
  "iconType": "stateFlower",
  "icon": "Black Eyed Susan"
}
```

The result of POST or PUT requests should use the following HTTP response codes:

* `200 OK` may be used if inserting with a `PUT` to differentiate between a new insertion and an update
* `201 CREATED` is used to state that a new resource was created in the collection. It is recommended that the ID of the new resource is returned in the response body, if applicable.
* `400 BAD REQUEST` should be used in the event of malformed JSON or an invalid request body
* `401 UNAUTHORIZED` should be used if we expect an authorization token to be provided on an endpoint and it isn't (see the [authentication docs](Authentication.md) for more info)
* `403 FORBIDDEN` should be used if an authorization token was provided and the requester doesn't have the appropriate roles or permissions to access the collection. The endpoint should explain what's missing in the response.
* `409 CONFLICT` is used if this collection has certain constraints on the contained data, this may indicate insertion of a duplicate resource or some other conflict with existing data

### Collection attributes

Collections may have attributes, such as the "state-icons" attribute in the second example in the [adding to the collection](#adding-to-the-collection)
section. These attributes should be preferred over RPC requests on collections, where possible.

For example, if you wanted the average population of all countries you might add an attribute-like path 
to your countries API:

`GET /api/v1/countries/population/average`

This route now has a "population" attribute for the collection of countries, which itself has an "average" attribute. 
Make sure to use adjectives and nouns when making these properties, not verbs.

Do this: `GET /api/v1/countries/population/average`

Don't do this: `GET /api/v1/countries/calculate-average-population`

Or this: `GET /api/v1/countries/population/calculate-average`

## Unique Resources

Unique resources within a collection can be retrieved via a `GET` request for their unique identifier. If you have
specific resources you expect to be in the system, there can be dedicated routes for them. In that case, use an alias
for the specific resource. Here are some examples:

* `GET /api/v1/countries/1` - Get the country with an ID of 1
* `GET /api/v1/countries/usa` - We expect the USA to be in the list of countries, so we have a special named reference for it in the collection
* `GET /api/v1/users/1` - Get the user with the ID 1
* `GET /api/v1/users/admin` - Specifically, get the admin user because we expect that to be in the system

The following response codes should be used when looking for unique resources:

* `200 OK` should be returned if the specified resource is found within the collection
* `401 UNAUTHORIZED` should be used if we expect an authorization token to be provided on an endpoint and it isn't (see the [authentication docs](Authentication.md) for more info)
* `403 FORBIDDEN` should be used if the requested resource is present but the user doesn't have the right roles or permissions to access the requested resource. The response body should explain what's missing.
* `404 NOT FOUND` should be returned if the specified resource is not found. If looking for nested resources, the response body should indicate which part was missing.
  * For example, if looking up `GET /api/v1/countries/4/provinces/1` returns a `404 NOT FOUND` because country 4 doesn't exist, the response body should state country 4 didn't exist.

### Updating resources

Resources should be updated by ID via the HTTP verbs `PUT` or `PATCH`. There's a slight difference between the two, though:

* `PUT` is an update that fully replaces the requested resource
* `PATCH` is a partial update that just specifies the couple of fields in the resource to change, leaving everything else alone

Here's an example of a full update:

```http request
PUT /api/v1/countries/5
Content-Type: application/json

{
  "name": "Genovia the Great",
  "flagUrl": "https://en.wikipedia.org/wiki/File:Flag_of_Genovia.svg"
}
```

Here's an example of a partial update:

```http request
PATCH /api/v1/countries/5
Content-Type: application/json

{
  "name": "Genovia the Great"
}
```

In both cases, the following response codes should be used:

* `200 OK` is used for a successful update
* `400 BAD REQUEST` is used for an invalid or malformed JSON payload
* `401 UNAUTHORIZED` should be used if we expect an authorization token to be provided on an endpoint and it isn't (see the [authentication docs](Authentication.md) for more info)
* `403 FORBIDDEN` should be used if the requested resource is present but the user doesn't have the right roles or permissions to access or update the requested resource. The response body should explain what's missing.
* `404 NOT FOUND` is used when trying to update a resource that doesn't exist. If updating a nested resource, the missing resource should be specified in the error message. This does not apply if using `PUT` with "upsert" semantics.
* `409 CONFLICT` is used for an update that would conflict with existing data, such as if a certain field needed to be unique among the other members of the collection

### Deleting resources

You may delete a resource by using the `DELETE` verb targeting a resource. Deleting Genovia in the previous examples
might look like this: `DELETE /api/v1/countries/5`.

The following HTTP responses should be used when deleting resources:

* `200 OK` is used for a successful deletion
* `401 UNAUTHORIZED`should be used if we expect an authorization token to be provided on an endpoint and it isn't (see the [authentication docs](Authentication.md) for more info)
* `403 FORBIDDEN` should be used if the requested resource is present but the user doesn't have the right roles or permissions to access or delete the requested resource. The response body should explain what's missing.
* `404 NOT FOUND` is used when trying to delete a resource that doesn't exist. If a parent resource above a nested resource doesn't exist, that should be specified in the error response.
* `409 CONFLICT` is used if the operation disrupts data integrity, such as if this resource is depended on by other resources which need to be deleted first

## Triggering actions on resources

In some cases, we may need to trigger some process in the backend which isn't necessarily a CRUD operation and doesn't
fit nicely into the previous set of guidelines for building out APIs. This is where "RPC-like" actions on resources come
into play.

"RPC-like" actions on resources are triggered with HTTP `POST` verbs. They should not accept query parameters and any
inputs necessary for the action should be provided in the request body. This is the only case where you may use verbs
in URLs - verbs signify the endpoint is a resource action rather than a subresource.

Here's an example of a resource action:

```http request
POST /api/v1/countries/usa/states/maryland/deliver-mail
Content-Type: application/json

{
  "targetZipCodes": [21044, 21202]
}
```

The following HTTP responses should be used for resource actions:

* `200 OK` should indicate the operation completed successfully
* `400 BAD REQUEST` should be used if the request body is malformed or has invalid values
* `401 UNAUTHORIZED`should be used if we expect an authorization token to be provided on an endpoint and it isn't (see the [authentication docs](Authentication.md) for more info)
* `403 FORBIDDEN` should be used if the requested resource is present but the user doesn't have the right roles or permissions to perform the action. The response body should explain what's missing.
* `404 NOT FOUND` should be used if the resource we're performing the action on or a parent resource does not exist. If a parent resource does not exist, the response body should specify the parent resource was missing.

Any problems that occur during the action should be reported in a `500 INTERNAL SERVER ERROR`.

## Maintaining compatibility with old APIs when porting routes

If we are creating an endpoint that replaces an endpoint on another existing microservice (such as with a microservice migration)
and the old endpoint doesn't follow API naming conventions, that API endpoint may be added for compatibility reasons. 

As soon as the badly named endpoint is added, a ticket should be created for replacing it with a correctly named endpoint.
After the succeeding endpoint has been added, a deprecation cycle for the old endpoint should be followed as described 
in the [versioning section](#casing-and-versioning).

