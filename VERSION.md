git add . &&
git commit --amend --no-edit && 
git push origin staging -f &&

git tag v1.0.16 &&
git push origin v1.0.16