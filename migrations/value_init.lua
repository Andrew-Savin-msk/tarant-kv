box.cfg{}

box.schema.space.create('example_space', {if_not_exists = true})
box.space.example_space:format({
    {name = 'id', type = 'unsigned'},
    {name = 'value', type = 'string'},
})
box.space.example_space:create_index('primary', {
    type = 'hash',
    parts = {'id'},
    if_not_exists = true,
})
