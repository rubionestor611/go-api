# Nestor's Golang API
I'm taking a venture into learning some Golang for API development and plan for this one api to host several minor tools and cool ventures to come. The functionaliteis of the endpoints are listed below.

## Path Breakdown
| Path | Description |
| ----------- | ----------- |
| /v1/healthcheck | Endpoint defined in order to know whether an instance is healthy or not. A status of 200 means the instance is healthy. Anything else means it is unhealthy. |
| /nestor/spotify/currently-playing | Returns the song I am currently listening to if I am online. If not, a message saying I am not listening to a song is returned. This excludes podcasts as I want to flex my music taste. |
| /nestor/spotify/top-tracks | Returns my top 20 songs over a period of the last 4 weeks or so according to Spotify. |
| /nestor/spotify/profile | Returns information about my Spotify profile like my username, profile picture, follower count, and URL to my account. |
| /nestor/spotify/recently-played | Returns my most recently played song. This excludes podcasts as I felt I wanted to share my music taste more than anything. |
