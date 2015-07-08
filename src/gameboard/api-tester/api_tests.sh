#!/bin/bash
echo "start"
curl -X POST http://localhost:8000/game -d game_id=1
echo ""
curl -X GET "http://localhost:8000/game?game_id=1"
echo ""
curl -H "Content-Type: application/octet-stream" -X POST -d @shapedata.json "http://localhost:8000/shape?game_id=1" 
echo ""
echo "done"
