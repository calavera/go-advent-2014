package main

import (
	"os"
	"time"

	"github.com/libgit2/git2go"
)

func credentialsCallback(url string, username string, allowedTypes git.CredType) (git.ErrorCode, *git.Cred) {
	ret, cred := git.NewCredSshKeyFromAgent(username)
	return git.ErrorCode(ret), &cred
}

func certificateCheckCallback(cert *git.Certificate, valid bool, hostname string) git.ErrorCode {
	if hostname != "github.com" {
		return git.ErrUser
	}
	return 0
}

func main() {
	repo, err := git.Clone("git://github.com/gopheracademy/gopheracademy-web.git", "web", &git.CloneOptions{})
	if err != nil {
		panic(err)
	}

	signature := &git.Signature{
		Name:  "David Calavera",
		Email: "david.calavera@gmail.com",
		When:  time.Now(),
	}

	head, err := repo.Head()
	if err != nil {
		panic(err)
	}

	headCommit, err := repo.LookupCommit(head.Target())
	if err != nil {
		panic(err)
	}

	branch, err := repo.CreateBranch("git2go-tutorial", headCommit, false, signature, "Branch for git2go's tutorial")
	if err != nil {
		panic(err)
	}

	idx, err := repo.Index()
	if err != nil {
		panic(err)
	}

	err = os.Link("git2go-tutorial.md", "web/content/advent-2014/git2go-tutorial.md")
	if err != nil {
		panic(err)
	}

	err = idx.AddByPath("content/advent-2014/git2go-tutorial.md")
	if err != nil {
		panic(err)
	}

	treeId, err := idx.WriteTree()
	if err != nil {
		panic(err)
	}

	tree, err := repo.LookupTree(treeId)
	if err != nil {
		panic(err)
	}

	commitTarget, err := repo.LookupCommit(branch.Target())
	if err != nil {
		panic(err)
	}

	message := "Add Git2go tutorial"
	_, err = repo.CreateCommit("refs/heads/git2go-tutorial", signature, signature, message, tree, commitTarget)
	if err != nil {
		panic(err)
	}

	fork, err := repo.CreateRemote("calavera", "git@github.com:calavera/gopheracademy-web.git")

	cbs := &git.RemoteCallbacks{
		CredentialsCallback:      credentialsCallback,
		CertificateCheckCallback: certificateCheckCallback,
	}

	err = fork.SetCallbacks(cbs)
	if err != nil {
		panic(err)
	}

	push, err := fork.NewPush()
	if err != nil {
		panic(err)
	}

	err = push.AddRefspec("refs/heads/git2go-tutorial")
	if err != nil {
		panic(err)
	}

	err = push.Finish()
	if err != nil {
		panic(err)
	}
}
