package errgroup

import(
	"log"
	"golang.org/x/sync/errgroup"
)

type Result struct {
	userID int
	username string
	post *Post
	comment *Comment
}

func errGroupExample() {		
	//Retrieve this data with userID 1
	result := &Result{userID:1}
	
	//Using Errgroup
	g := errgroup.Group{}
	
	func(r *Result) {
		g.Go(func() error {
			var err error
			r.username, err = getUsername(r.userID)
			if err != nil {
				return err
			}
			return nil
		}) //Username
		
		g.Go(func() error {
			var err error
			r.post, err = getLastPost(r.userID)
			if err != nil {
				return err
			}
			return nil
		}) //Last Post
		
		g.Go(func() error {
			var err error
			r.comment, err = getLastComment(r.userID)
			if err != nil {
				return err
			}
			return nil
		}) //Last Comment
	}(result)
	
	//Wait for errgroup to finish
	err := g.Wait()
	if err != nil {
		log.Println(err)
	}
	
	//Show off retrieved data
	log.Println(result.username)
	log.Println(result.post)
	log.Println(result.comment)
	
	//Close database connection
	DB.Close()
}