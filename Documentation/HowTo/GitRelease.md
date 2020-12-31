## Git release on travis

Since tags are treated as branches in travis you have to use regex expressions to match tags:
```yaml
branches:
  only:
  - master
  - /v\d+\.\d+\.\d+/
```

The tag can be checked with an if condition for each stage:
```yaml
- stage: test
    name: "test tag"
    if: tag IS NOT present
```

The github release can be added to a stage where the tag is present:
```yaml
deploy:
  provider: releases
  api_key: $GITHUB_OAUTH_TOKEN
  file_glob: true
  file: ./out/make/**/*.exe
  skip_cleanup: true
  on:
    tags: true
```

Make sure to skip cleanup as this will remove all files you just created and the deployment will fail.

### Trigger release

First add all changes you made to your commit and push to origin.
```bash
git add .
git commit -m "test"
git push
```
This will trigger a build, but since we specified only certain stages will be run if a tag is present, only half of our specified stages will run. Now you can trigger another build with:
```
git tag v0.0.1
git push --tags
```
Pushing tags will trigger another build, this time all the other missing stages will run.

Releases can then be viewed by visiting: https://github.com/AndreasRoither/Snippets/releases