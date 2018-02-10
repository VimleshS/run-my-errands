# run-my-errands
The service app instantly connects you with trusted erranders to do your grocery shopping.  So you can be more productive every day. There are two services running behind this to make process seamless. Orderer is the main web application which places the order and dispatcher dispatches all the assembled grocery items and send email

#How this service works

### Request to generate authenticate and generate Jwt

   Authorize youself with email and password, this will generate a token send it back to server in authorization header for authenticated api call

    curl -XPOST 'https://runaerrander.herokuapp.com/authenticate' -H 'Content-Type: application/json' -d '
    {
        "email": "john.doe@gmail.com",
        "password": "mysecretpassword"
    }
    '

   Response 
    
    {"token":"dfdfdf"}
    

### Upload your grocery list

    curl -XPOST 'https://runaerrander.herokuapp.com/uploadlist' \
    -H 'authorization: bearer dfdfdf' \
    -H 'Content-Type: application/json' -d '
    {
        "id": 1,
        "total": 110,
        "groceries": [{
            "name": "Item1",
            "quantity": 50
        }, {
            "name": "Item2",
            "quantity": 60
        }]
    }
    '

   A sample response 
    
    {
    "id": 1,
    "email": "john.doe@gmail.com",
    "total": 110,
    "groceries": [{
      "name": "Item1",
      "quantity": 50
    }, {
      "name": "Item2",
      "quantity": 60
    }],
    "message": "Order processed and forwarded for dispatching"
    }



#### You will receive a dispatch notification on your email.
