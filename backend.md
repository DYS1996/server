# Backend

## Logical
Posts // article posted
- title // unique, unnullable
- created date // unnullable
- last modified date
- tags // multi, all different for same post(unique(postID, tag))
- content // unnullable 
- **postID**

Comments // comments to article
- authorEmail // unnullable
- created date // unnullable
- content // unnullable
- postID // unnullable, fk to Posts(postID)
- **commentID**

Users // users of this blog site
- **uid**
- userName // unnullable unique
- passWord // unnullable sha256 
- privilege // unnullable default 100

## Physical

### Tables:

Posts
- title // index, text, unnullable, length: [3,15], unique
- cDate // date, unnullable, default current_date
- mDate // date, default null 
- content // text, unnullable, length: [10, 3500]
- postID // SERIAL, pk
- fullTextSearch // tsvector, computed on (title, content), patch-1
index(fullTextSearch)

Tags
- tagID // SERIAL, pk
- postID // int, fk -> Posts(postID), unnullable
- tag // text, unnullable, length: [2, 6]
constraints: unique(postID, tag), one post has no more than 5(tag)

Comments
- commentID // SERIAL, pk
- postID // int, fk -> Posts(postID), unnullable
- cDate // date, unnullable, default current_date
- authorEmail // text, unnullable
- content // text, unnullable, length:[2, 100]
index(postID)

Users // patch-2
- uid // SERIAL, pk
- userName // text, unique, unnullable, len: [5, 14]
- passWord // bytea, unnullable 
- privilege // int, unnullable, default 100, constraint: [0,100]

### ViewTypes:

PostView
- postID INT
- title TEXT
- cDate DATE
- mDate DATE
- content TEXT
- tags TEXT[]

CommentView 
- postID INT
- commentID INT
- authorEmail TEXT
- cDate DATE
- content TEXT

UserView // patch-2
- uid INT
- userName TEXT
- privilege INT

### APIs:

getPostByID(pid INT): setod PostView

getPostsByPage(pagesize INT, page INT): setof PostView

getPostsCount(): INT

insertPost(title TEXT, content TEXT, tags TEXT[]): INT

deletePost(pid INT) BOOLEAN

updatePost(pid INT, newTitle TEXT, newContent TEXT, newTags TEXT[]): BOOLEAN

getCommentsCount(pid INT): INT

getCommentsByPage(pid INT, pagesize INT, page INT): setof CommentView

insertComment(pid INT, content TEXT, authorEmail TEXT): INT

deleteComment(cmtID INT): BOOLEAN

updateComment(cmtID INT, newContent TEXT, newAuthorEmail TEXT): BOOLEAN

getPostsByFTS(query TEXT, page INT, pagesize INT): setof PostView // patch-1

getPostsCountByFTS(query TEXT): INT // patch-1

getUser(user_name TEXT, pass BYTEA): setof UserView // patch-3

insertUser(user_name TEXT, pass BYTEA): INT // patch-3

updateUser(userID INT, nPW BYTEA): BOOLEAN // patch-3