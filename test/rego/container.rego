package container

default allow = false

allow {
    gbac.can("a", "admin", "widget", "sub")
}

test_allow {
    allow
}