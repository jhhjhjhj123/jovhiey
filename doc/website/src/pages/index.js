import React from "react";
import LayoutProvider from "@theme/Layout/Provider";
import { AtlasGoWebsite } from "@ariga/atlas-website";

import "@ariga/atlas-website/style.css";

export default function () {
  return (
    <LayoutProvider>
      <AtlasGoWebsite
        events={{
          // TODO: Guys please, add your events here :)
          events: [],
          isHidden: true,
        }}
      />
    </LayoutProvider>
  );
}
