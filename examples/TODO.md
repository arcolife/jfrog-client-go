- [ ] repo creation from config params in yaml
- [ ] replace fmt.Print* with logging module (fix)
- [ ] activate and override Files param over Pattern (or option to include both)

- [x] Multiple file with pattern indeterministic
      we get 400 (Bad Request) or 201 (Accepted OK) randomly
      -> we can only use Threads=1 with free bintray version *apparently*

- [x] fix Path in UploadParams struct to avoid incomplete package paths in messages / eliminate untested APIs / fix potential bugs
      example incomplete path
      => "message": "Calculation was successfully scheduled for '/arcolife/rpm'"
      should be -> arcolife/rpm/CasperLabs
