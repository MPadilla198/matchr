# Contents

## Sources
Images are self provided to act as testing samples. [Sewar](https://pypi.org/project/sewar/) Python package was used to generate resulting calculations for simple image metrics.

## Commands Used
pip install sewar
for metr in mse rmse psnr rmse_sw uqi ssim ergas scc rase sam msssim vifp psnrb; do for gt in *.jpg; do for p in *.jpg; do echo "$gt $p $(sewar $metr $gt $p)"; done; done; done | tee results

