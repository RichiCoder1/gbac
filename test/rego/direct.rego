package direct

default allow = false

allow {
    gbac.can("d", "write", "widget:a", "default")
}

test_allow {
    allow
}