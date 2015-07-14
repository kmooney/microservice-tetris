#!/bin/bash
curl -X POST http://localhost:8000/game -d game_id=1
curl -H "Content-Type: application/octet-stream" -X PUT -d @shapedata.json "http://localhost:8000/shape?game_id=1" 
