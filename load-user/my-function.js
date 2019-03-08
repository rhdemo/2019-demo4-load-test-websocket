module.exports = { gameMock };
var nanoid = require('nanoid')
var jsonfile = require('jsonfile')
const uuidv1 = require('uuid/v1');
const file = 'sample.json'
var gesture;


function gameMock(userContext, events, done) {
  gesture = jsonfile.readFileSync(file)

  gesture['playerID'] = uuidv1()
  gesture['uuid'] = nanoid()


  // set the "data" variable for the virtual user to use in the subsequent action
  userContext.vars.data = gesture;
  return done();
}
