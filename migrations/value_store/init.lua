box.cfg{}

example_space = box.schema.create_space('values', {if_not_exists = true, format = {
    {name = 'key', type = 'string'},
    {name = 'value', type = 'scalar'},
}})
example_space:create_index('primary', {
    type = 'hash',
    parts = {'key'},
    if_not_exists = true,
})
