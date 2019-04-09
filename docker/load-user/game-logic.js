module.exports = { gameMock };
var nanoid = require('nanoid')
var jsonfile = require('jsonfile')
const uuidv1 = require('uuid/v1');
const file = 'sample.json'
var gesture;

gesture = jsonfile.readFileSync(file)


function gameMock(userContext, events, done) {

    userContext.vars.data = Object.assign({}, gesture, { uuid: nanoid() })
    done();
}
