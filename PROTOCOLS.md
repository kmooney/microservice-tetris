# Protocols

This document describes the protocols that each 
microservice will use to communicate with other 
microservices.  Each service will be built
so that it takes a request and provides a response.
The game, by default, is paused, and only becomes
unpaused when a `tick` message is sent to the 
gameboard.

## Gameboard

### Accepts

#### tick

The gameboard expects the tick message to be sent 
by the client.  We can expect between 2 to 10 ticks
per client * second.

##### Properties

* game id

#### game

##### POST

The gameboard expects new games to be created by
clients.  A client can have multiple games.

##### Properties:

* client_id

##### Response

When POSTed to, with a client ID, the gameboard
service creates a new gameboard and returns a
response that contains a status message.
The status message will have the form of 
`game_status` below.

##### GET

When a GET request is sent and a `game_id` is also provided
the game board responds with a `game_status`, detailed below.

#### subscription

##### POST 

When a subscription is created, it is a client telling the
gameboard service that when the status of the game changes
(for instance if the score changes, the placed shapes
or orientations change, or the line count changes) the
gameboard will create an HTTP request, either a GET or POST, 
depending on the details of the subscription, and it will 
render the template stored.

This allows clients to request that changes be pushed to them
rather than constantly polling the gameboard service.  After all, 
the gameboard service will already be handling 2-10 ticks per 
client second.

###### Properties

* game_id
* property
* response_method (GET or POST)
* response_url

##### Response

It will either respond with a 200 OK or an error message, see
`response` below.

#### validate_shape

##### POST

When a shape is submitted, it is checked for whether it is currently
in a valid position.  If the position is valid, returns 200 OK, if not
valid, returns HTTP 412 - Precondition failed (the precondition being that
the shape is in a valid position on the board).

###### Properties

* game_id
* shapedata

#### place_shape

##### POST

When a shape is submitted, if it can be placed, the shapedata for the
given gameboard will be updated to reflect that this shape has been
placed.  

If the shape is placed, the http response will be 200 (OK), if it cannot
be placed the response will be 412 - Precondition Failed (in this case, 
the precondition is that the shape is placeable)

It could be possible for the client to independently track whether the 
shape can be placed and only call this function when it believes that the
shape is placeable.

###### Properties

* game_id
* shapedata

##### Callback

It will send a `game_status` message as detailed below, with the
granularity of the `property` of the status requested.  If the
`property` field is left blank in the request, the entire game
status will be sent in the callback.

### Sends

#### response
     
The `response` message is always sent as a response to 
any HTTP request.  The response always contains an error,
a message, or both.

##### Properties

* Error - the error message and code.
* Message - the response message.  Any of the message types below.

#### game_status 

The `game_status` contains the current status of the game, including
shapes that have been placed, and their positions and orientations, 
the current line count, the current score and the current level.

This is the master dump of the status of the requested game.  
Users can also request more discrete information by drilling into the status
in the request.
 
For instance, one could request `/game_status/<<game_id>>/` for 
the entire status, but one could also request `/game_status/<<game_id>>/score` for 
just the score.  

##### Properties

* lines - the line count
* score - the current score
* board - a list of shapes, with position and orientation
* level - the current level, also used by the client to determine the tick speed.

## Rotator

The rotator service takes a serialized shape object,
rotates either clockwise or counterclockwise depending on 
the request, and returns a new shape object (which should in turn
replace the current shape and position as maintained by the _Current Shape_
service, below).

### Accepts

#### rotate

##### POST

##### Shapedata - A list of shape information:
* shapewidth 
* shapepresence 

* Direction - "CCW" or "CW"

##### Response

##### Response with the following message (see Response, above)
    * Shapedata

## Shape Generator

### Accepts

#### Shaperequest

##### GET
Response with the following message (see Response, above)
* Shapedata

## Current Shape Service

Current shape maintains some state.  It maintains the current
(x,y) cooridinates, shape data, and rotation.

The current shape service depends on the gameboard.
It has to know whether the piece has hit the bottom of the board, 
and whether the Current Shape PUT is valid.

### Accepts

#### Tick
##### GET
Expects a request with the following:
* Game ID

#### Current Shape
##### GET
Expects a request with the following:
* Game ID

Responds with the following message (see Response, above)
* Shapedata

##### PUT
Expects a message with the following contents:
* Game ID
* Shapedata

Responds with either a 200 OK to indicate that the submitted
Shapedata (with optional position data) is valid, or it responds
with an error that indicates that the position data indicated in the PUT is
invalid and that the Current Shape has not been updated to that.

### Sends
### Callbacks
#### Shape Change Event
Clients can register for a shape change event, by sending a Callback request.
When the callback is executed, a Shapedata is sent to the client who registered for the
request.
