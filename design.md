# Middleware

## Interface

/post
- GET: ?id: int --queryPost--> {err: null, data: {pid: int, title: string, cData: dateString, mDate: dateString, content: string, tags: [string]}}
- POST: auth need
    - {action: "insert", title: string, content: string, tags: [string]} --insertPost--> {err: null, data(pid): int}
    - {action: "delete", pid: int} --deletePost--> {err: null, data(pid): -1}
    - {action: "update", pid: int} --updatePost--> {err: null, data(pid): -1}

/posts
- GET: ?[keyword: string &] page: int & pageSize: int --queryPosts--> {err: null, data: {maxPage: int, posts: [{pid: int, title: string, cDate: dateString, mDate: dateString, content: string, tags: [string]}}]}

/comments
- GET: ?pid: int & page: int & pageSize: int --queryPost--> {err: null, data: {maxPage: int, comments: [{pid: int, cid: int, email: emailString, cDate: dateString, content: string}]}}

/comment
- POST: auth need
    - {action: "insert", pid: int, content: string, authorEmail: string} --insertComment--> {err: null, data(cid): int}
    - {action: "delete", commentID: int} --deleteComment--> {err: null, data(cid): -1}
    - {action: "update", commentID: int, newContent: string, newEmail: emailString} --updateComment--> {err: null, data(cid): -1}

/user
- POST
    - {action: "login", userName: string, passWord: string} --loginUser--> {err: null, data: {uid: int, userName: string, privilege: int}}
    - {action: "register", userName: string, passWord: string} --insertUser--> {err: null, data(uid): int}
    - {action: "update", uid: int, newPassWord: string} --updateuser--> {err:null, data(uid): -1}

/ping
- --pingTest--> "pong"