package general

default allow = false

allow {
    gbac.can("a", "admin", "widget:a", "default")
}

test_allow {
    allow
}