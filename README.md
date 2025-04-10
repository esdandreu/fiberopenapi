# fiberopenapi
Fiber router interface generation using openapi specifications

Spec first!

Inspired by https://github.com/openapistack/openapi-backend.

Based on code generators instead.


TODO:
- Split in files but use a generator to put them together in a single file.

## Why should you use this?

You should only use this if you plan to define your API with an OpenAPI
specification.

Generated code has some advantages compared to other frameworks:

- No obscure middleware. Generated code goes through your code analysis, passes
  your code review and satisfies your quality requirements.
- No copyright headaches. Third party code is not used in your final binary.
  This is a third party tool. Small enough that you can simply copy it as a
  script into your project.
- No unused logic. Generate only the features you are going to use. Modify the
  generator to fit your needs.