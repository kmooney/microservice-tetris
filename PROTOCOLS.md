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

    Properties:

    * game id

#### game

    _POST_

    The gameboard expects new games to be created by
    clients.  A client can have multiple games.

    Properties:

    * client_id

    Response:

    When POSTed to, with a client ID, the gameboard
    service creates a new gameboard and returns a
    response that contains a status message.
    The status message will have the form of 
    `game_status` below.

    _GET_

    When a GET request is sent and a `game_id` is also provided
    the game board responds with a `game_status`, detailed below.

#### subscription

    _POST_ 

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

    Properties

    * game_id
    * property
    * response_method (GET or POST)
    * response_url

    ##### Response

    It will either respond with a 200 OK or an error message, see
    `response` below.

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

     Properties:
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

     *Properties*
         * lines - the line count
         * score - the current score
         * board - a list of shapes, with position and orientation
         * level - the current level, also used by the client to determine
         the tick speed.

