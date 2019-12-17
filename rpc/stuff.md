

# Login/Signup
{"/user/register", "POST", "name", "username", "pwhash", "seedpwd"} for registering as a user on the main screen  
{"/token", "POST"} to generate a token
{"/user/validate", "GET"} for logging onto the platform, getting details about a specific user, any info needed on investors, etc  

# Register / Validate Entities
{"/investor/register", "POST", "name", "username", "pwhash", "token", "seedpwd"} to register new entities
{"/recipient/register", "POST", "name", "username", "pwhash", "seedpwd"} to register new entities
{"/entity/register", "POST", "name", "username", "pwhash", "token", "seedpwd", "entityType"} to register new entities

{"/investor/validate", "GET"} used to register as particular entities and getting data about them  
{"/recipient/validate", "GET"} used to register as particular entities and getting data about them  
{"/entity/validate", "GET"} used to register as particular entities and getting data about them                                                                   
# Dashboard Endpoints
{"/investor/dashboard", "GET"} to load partial information on the dashboard screens  
{"/recipient/dashboard", "GET"} to load partial information on the dashboard screens  
{"/entity/contractor/dashboard", "GET"} to load partial information on the dashboard screens  
{"/entity/developer/dashboard", "GET"} to load partial information on the dashboard screens  

# Ipfs endpoints
{"/ipfs/getdata", "GET", "hash"} to retrieve data from ipfs  
{"/ipfs/putdata", "POST", "data"} to store data in ipfs  

# Stellar Endpoints
{"/user/latestblockhash", "GET"} to get the latest blockhash from the Stellar blockchain  
{"/user/sendxlm", "GET", "destination", "amount", "seedpwd"} to send xlm from the platform  
{"/user/trustasset", "GET", "assetCode", "assetIssuer", "limit", "seedpwd"} to trust a p2p asset  
{"/user/balances", "GET"} to load all user balances  
{"/user/balance/xlm", "GET"} to load the XLM balance  
{"/user/balance/asset", "GET", "asset"} to load the asset balance associated with the platform  

# Recovery shares
{"/user/sendrecovery", "GET", "email1", "email2", "email3"} to send recovery shares to email entities  
{"/user/seedrecovery", "GET", "secret1", "secret2"} to recover your lost seed  
{"/user/newsecrets", "GET", "seedpwd", "email1", "email2", "email3"} to generate new shares  

# Reset password
{"/user/resetpwd", "GET", "seedpwd", "email"} to send the login code via email  
{"/user/pwdreset", "GET", "pwhash", "email", "verificationCode"} to actually set a new password  

# Two factor authentication
{"/user/2fa/generate", "GET"} to generate the 2FA secret  
{"/user/2fa/authenticate", "GET", "password"} to authenticate with an external app (like Google Authenticator)  

# Invest on the platform
{"/investor/invest", "POST", "seedpwd", "projIndex", "amount"} to invest  
{"/recipient/unlock/opensolar", "POST", "seedpwd", "projIndex"} to unlock an investment  

# Update User params
{"/user/update", "POST"} name, city, pwhash, zipcode, country, recoveryphone, address, description, email, notification  

# Sweep funds
{"/user/sweep", "GET", "seedpwd", "destination"} to sweep funds from an account  
{"/user/sweepasset", "GET", "seedpwd", "destination", "assetName", "issuerPubkey"} to sweep specific assets from an address to another address  

# Teller related stuff
{"/tellerping", "GET"} to check if the teller is up  
{"/recipient/teller/details", "POST", "projIndex", "url", "brokerurl", "topic"} get the details of the teller  