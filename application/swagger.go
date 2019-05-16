/*
Copyright 2019 Adobe. All rights reserved.
This file is licensed to you under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License. You may obtain a copy
of the License at http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software distributed under
the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR REPRESENTATIONS
OF ANY KIND, either express or implied. See the License for the specific language
governing permissions and limitations under the License.
*/

package application

import (
	"github.com/kataras/iris"
)

const index = `
<!-- HTML for static distribution bundle build -->
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <title>OSSIA Reference</title>
    <link rel="stylesheet" type="text/css" href="assets/swagger-ui.css" >
    <link rel="icon" type="image/png" href="assets/favicon-32x32.png" sizes="32x32" />
    <link rel="icon" type="image/png" href="assets/favicon-16x16.png" sizes="16x16" />
    <style>
      html
      {
        box-sizing: border-box;
        overflow: -moz-scrollbars-vertical;
        overflow-y: scroll;
      }

      *,
      *:before,
      *:after
      {
        box-sizing: inherit;
      }

      body
      {
        margin:0;
        background: #fafafa;
      }
      /* section.models {
        display: none;
      } */
    </style>
  </head>

  <body>
    <div id="swagger-ui"></div>

    <script src="assets/swagger-ui-bundle.js"> </script>
    <script src="assets/swagger-ui-standalone-preset.js"> </script>
    <script>
    window.onload = function() {
      //document.getElementsByClassName("version").innerHTML = "new content"
      // Build a system
      const ui = SwaggerUIBundle({
        url: "/assets/swagger.json",
        dom_id: '#swagger-ui',
        deepLinking: true,
        defaultModelsExpandDepth: -1,
        presets: [
          SwaggerUIBundle.presets.apis,
          //SwaggerUIStandalonePreset.slice(1)
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        //layout: "StandaloneLayout"
      })

      window.ui = ui
    }
  </script>
  </body>
</html>
`

func apiReference(c iris.Context) {
	c.HTML(index)
}
