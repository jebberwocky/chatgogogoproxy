info "bfg sensitive data"
bfg --replace-text replacements.txt
git reflog expire --expire=now --all && git gc --prune=now --aggressive
git push origin main --force