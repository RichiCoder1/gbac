package inherited

default allow = false

allow {
    gbac.can("a", "admin", "widget:a", "default")
    gbac.can("a", "admin", "widget", "default")
    gbac.can("a", "write", "widget", "default")
    gbac.can("a", "read", "widget", "default")
}

test_allow {
    allow
}