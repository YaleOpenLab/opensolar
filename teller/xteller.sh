# List of endpoints to go through with Tornike
# Register User
curl -X POST -H "Content-Type: application/x-www-form-urlencoded" -d 'username=torniketest&pwhash=e9a75486736a550af4fea861e2378305c4a555a05094dee1dca2f68afea49cc3a50e8de6ea131ea521311f4d6fb054a146e8282f8e35ff2e6368c1a62e909716&seedpwd=x&email=test4@test5.com' "http://localhost:8081/user/register"

# Get a token
curl -X POST -H "Content-Type: application/x-www-form-urlencoded" -d 'username=torniketest&pwhash=e9a75486736a550af4fea861e2378305c4a555a05094dee1dca2f68afea49cc3a50e8de6ea131ea521311f4d6fb054a146e8282f8e35ff2e6368c1a62e909716' "http://localhost:8081/token"

# Login using the /validate endpoint
curl -X GET -H "Content-Type: application/x-www-form-urlencoded" -H "Origin: localhost" "http://localhost:8081/user/validate?username=torniketest&token=bCPUeuycOqjZARvhioHrCjUacAbykNdB"

# Register as an investor
curl -X POST -H "Content-Type: application/x-www-form-urlencoded" -d 'username=torniketest&pwhash=e9a75486736a550af4fea861e2378305c4a555a05094dee1dca2f68afea49cc3a50e8de6ea131ea521311f4d6fb054a146e8282f8e35ff2e6368c1a62e909716&token=bCPUeuycOqjZARvhioHrCjUacAbykNdB&name=TornikeTest&seedpwd=x' "http://localhost:8081/investor/register"

# Validate Investor
curl -X GET "http://localhost:8081/investor/validate?username=torniketest&token=bCPUeuycOqjZARvhioHrCjUacAbykNdB"

# View Investor Dashboard with some fields filled (can fill mroe if desired)
curl -X GET "http://localhost:8081/investor/dashboard?username=torniketest&token=bCPUeuycOqjZARvhioHrCjUacAbykNdB"

# Get XLM to setup account
curl -X GET "http://localhost:8081/user/askxlm?username=torniketest&token=bCPUeuycOqjZARvhioHrCjUacAbykNdB"

# Get Stablecoin to invest in a project, use amount as 1 to get 10000000 USD
curl -X GET "http://localhost:8081/stablecoin/get?username=torniketest&token=bCPUeuycOqjZARvhioHrCjUacAbykNdB&seedpwd=x&amount=1"

# Invest in the project
curl -X POST -H "Content-Type: application/x-www-form-urlencoded" -d 'username=torniketest&token=bCPUeuycOqjZARvhioHrCjUacAbykNdB&seedpwd=x&projIndex=1&amount=10' "http://localhost:8081/investor/invest"