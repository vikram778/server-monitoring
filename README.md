# server-monitoring
Monitoring various services for their uptime.

Steps to execute.

1. Update the env variables in env.local file 
2. go build to build the binary 
3. execute the binary

should be able to see the reports getting generated in the path mentioned in the env.local

#Improvement scope

1. For maintaining the list of servers would be ideal to keep them in DB .
2. So once the list of the servers are in db at the time of execution all the server list can be pulled and put it in a cache.
