module.exports = { gameMock };
var nanoid = require('nanoid')
var jsonfile = require('jsonfile')
const uuidv1 = require('uuid/v1');
const file = 'reduce-sample.json'
var gesture;


function gameMock(userContext, events, done) {
  gesture = {}
  gesture['type'] = "load-test";
  gesture['machineId'] = Math.floor(Math.random() * Math.floor(10));
  gesture['sensorId'] = uuidv1();
  gesture['vibrationClass'] = "kafka-test";
  gesture['confidence'] = Math.floor(Math.random() * Math.floor(100));


  // set the "data" variable for the virtual user to use in the subsequent action
  userContext.vars.data = gesture;
  return done();
}
