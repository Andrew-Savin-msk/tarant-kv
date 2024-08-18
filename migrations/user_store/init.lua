box.cfg{}

example_space = box.schema.create_space('users', {if_not_exists = true, format = {
    {name = 'login', type = 'string'},
    {name = 'password', type = 'string'},
}})
example_space:create_index('primary', {
    type = 'hash',
    parts = {'login'},
    if_not_exists = true,
})
