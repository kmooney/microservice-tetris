#!/bin/bash
echo "Creating Game ID # 1"
echo ""
curl -X POST http://localhost:8000/game -d game_id=1
echo ""
curl -X POST -d "game_id=1&response_url=localhost:9999" "http://localhost:8000/subscribe" 
echo ""
curl "http://localhost:8000/subscribe?game_id=1"
echo ""
