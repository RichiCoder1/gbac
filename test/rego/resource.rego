package resource

default allow = false

allow {
    gbac.can("e", "read", "widget:a", "default")
}

test_allow {
    allow
}