"""Automated unit test verifying format v2 precipitation output formatting."""
import subprocess
import unittest


class TestV2Precipitation(unittest.TestCase):
    def test_v2_renderer_precipitation_metric_and_imperial(self):
        result = subprocess.run(
            ["go", "test", "-v", "./internal/renderer/v2/..."],
            capture_output=True,
            text=True,
        )
        self.assertEqual(
            result.returncode,
            0,
            f"go test failed:\nstdout:\n{result.stdout}\nstderr:\n{result.stderr}",
        )
        self.assertIn("PASS: TestDrawRainSparkMetricUnits", result.stdout)
        self.assertIn("PASS: TestDrawRainSparkImperialUnits", result.stdout)
        self.assertIn("PASS: TestV2Renderer_RenderPrecipitationMetric", result.stdout)
        self.assertIn("PASS: TestV2Renderer_RenderPrecipitationImperial", result.stdout)


if __name__ == "__main__":
    unittest.main()
