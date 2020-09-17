# OPA GBAC

This is a toy project that tries to model a persitent graph-based permissions structure and make it easy to reference from OPA. This is intended to be used along side traditional business logic. Eg:

```go
package widget

default allow = false

allow {
    input.action == "write"
    can_write
}

can_write {
    input.widget.owner = input.jwt.payload.sub
}

can_write {
    input.jwt.payload.claim["admin"]
}

can_write {
    gbac.can(input.jwt.payload.sub, "write", "widget:" + input.widget.id, "system")
}
```

This projects makes a couple of assumptions, all of which could be made more flexible with more iteration:
1) Users are an atomic unit, with no heirarchy
2) "Actions" are heirarchal, with some actions potentially allowing others
3) You want to have action permission for all resources (e.g. 'admin,write,read') or just specific resource types ('widget:write')
4) Resources are assiated with containers, aka things like accounts, which then might be in organizations. 
5) Granting an "action" permission to a parent org, gives that user permissions to child orgs
6) There are times you want to grant a user an "action" permission on a specific resource only

## Examples

To run examples, you'll first need a redisgraph instance. Easiest way to get a test one is `docker run -p 6379:6379 -it --rm redislabs/redisgraph`.

Then, you can add some test data via `gp run test/redis/internal/bootstrap.go`. You can examine that file to see the general graph structure of the test data.

Finally, you can `make run-opa` to start an instances of OPA with the `gpac` built-in enabled, and some test policies in `test/rego` loaded. 

Try, for example, running:
```sh 
curl -X POST -H "Content-Type: application/json" \
    -d '{}' \
    http://localhost:8181/v1/data/direct/allow
```

## `gbac.can`

This built-in takes four arguments:
* `user` - The identifier of the user performing the operation
* `action` - The action the `user` is attempting to perform
* `resource` - The `resource type` to perform, optionally with `resource id` in the form of `<type>:<id>`
* `owner` - The `container` containing the resource. This PoC assumes that most of the time, a `user` is granted permissions to resources in a `container` rather than directly.