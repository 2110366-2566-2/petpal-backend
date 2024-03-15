This generate
1. 5 User register
   - each has email : `[1..5]@user.com` and password : `password`
   - has own 5 pets with random attribute
2. 5 Service provider register
   - each has email : `[1..5]@svcp.com` and password : `password`
   - create 5 service
3. has booking
   - user index `i` books all the service of service provider index `i` with a **first** timeslot

**caution** in the booking process api will sent a lot of email (around 25 emails by default) so this will take much more longer if you generate more user or more service per svcp.