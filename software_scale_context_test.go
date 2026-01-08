package astiav

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSoftwareScaleContext(t *testing.T) {
	f1 := AllocFrame()
	require.NotNil(t, f1)
	defer f1.Free()

	f2 := AllocFrame()
	require.NotNil(t, f2)
	defer f2.Free()

	f3 := AllocFrame()
	require.NotNil(t, f3)
	defer f3.Free()

	srcW := 4
	srcH := 2
	srcPixelFormat := PixelFormatYuv420P
	dstW := 8
	dstH := 4
	dstPixelFormat := PixelFormatRgba
	swscf1 := SoftwareScaleContextFlags(SoftwareScaleContextFlagBilinear)

	f1.SetHeight(srcH)
	f1.SetWidth(srcW)
	f1.SetPixelFormat(srcPixelFormat)
	require.NoError(t, f1.AllocBuffer(1))

	swsc1, err := CreateSoftwareScaleContext(srcW, srcH, srcPixelFormat, dstW, dstH, dstPixelFormat, swscf1, nil, nil, []float64{1.1, 2.2})
	require.NoError(t, err)
	defer swsc1.Free()
	require.Equal(t, 1.1, swsc1.ScalerParam0())
	require.Equal(t, 2.2, swsc1.ScalerParam1())
	require.Equal(t, dstW, swsc1.DestinationWidth())
	require.Equal(t, dstH, swsc1.DestinationHeight())
	require.Equal(t, dstW, swsc1.DestinationWidth())
	require.Equal(t, swscf1, swsc1.Flags())
	require.Equal(t, srcH, swsc1.SourceHeight())
	require.Equal(t, srcPixelFormat, swsc1.SourcePixelFormat())
	require.Equal(t, srcW, swsc1.SourceWidth())
	require.Equal(t, srcH, swsc1.SourceHeight())
	require.Equal(t, srcW, swsc1.SourceWidth())
	cl := swsc1.Class()
	require.NotNil(t, cl)
	require.Equal(t, "SWScaler", cl.Name())

	require.NoError(t, swsc1.ScaleFrame(f1, f2))
	require.Equal(t, dstH, f2.Height())
	require.Equal(t, dstW, f2.Width())
	require.Equal(t, dstPixelFormat, f2.PixelFormat())

	dstW = 4
	dstH = 3
	dstPixelFormat = PixelFormatYuv420P
	swscf2 := SoftwareScaleContextFlags(SoftwareScaleContextFlagPoint)
	srcW = 2
	srcH = 1
	srcPixelFormat = PixelFormatRgba

	require.NoError(t, swsc1.SetDestinationHeight(dstH))
	require.Equal(t, dstH, swsc1.DestinationHeight())
	require.NoError(t, swsc1.SetDestinationPixelFormat(dstPixelFormat))
	require.Equal(t, dstPixelFormat, swsc1.DestinationPixelFormat())
	require.NoError(t, swsc1.SetDestinationWidth(dstW))
	require.Equal(t, dstW, swsc1.DestinationWidth())
	dstW = 5
	dstH = 4
	require.NoError(t, swsc1.SetDestinationWidth(dstW))
	require.NoError(t, swsc1.SetDestinationHeight(dstH))
	require.Equal(t, dstW, swsc1.DestinationWidth())
	require.Equal(t, dstH, swsc1.DestinationHeight())
	require.NoError(t, swsc1.SetFlags(swscf2))
	require.Equal(t, swsc1.Flags(), swscf2)
	require.NoError(t, swsc1.SetSourceHeight(srcH))
	require.Equal(t, srcH, swsc1.SourceHeight())
	require.NoError(t, swsc1.SetSourcePixelFormat(srcPixelFormat))
	require.Equal(t, srcPixelFormat, swsc1.SourcePixelFormat())
	require.NoError(t, swsc1.SetSourceWidth(srcW))
	require.Equal(t, srcW, swsc1.SourceWidth())
	srcW = 3
	srcH = 2
	require.NoError(t, swsc1.SetSourceWidth(srcW))
	require.NoError(t, swsc1.SetSourceHeight(srcH))
	require.Equal(t, srcW, swsc1.SourceWidth())
	require.Equal(t, srcH, swsc1.SourceHeight())

	require.NoError(t, swsc1.SetSourceHorizontalChromaPosition(1))
	require.Equal(t, 1, swsc1.SourceHorizontalChromaPosition())
	require.NoError(t, swsc1.SetSourceVerticalChromaPosition(2))
	require.Equal(t, 2, swsc1.SourceVerticalChromaPosition())
	require.NoError(t, swsc1.SetDestinationHorizontalChromaPosition(3))
	require.Equal(t, 3, swsc1.DestinationHorizontalChromaPosition())
	require.NoError(t, swsc1.SetDestinationVerticalChromaPosition(4))
	require.Equal(t, 4, swsc1.DestinationVerticalChromaPosition())
	require.NoError(t, swsc1.SetScalerParam0(5))
	require.Equal(t, float64(5), swsc1.ScalerParam0())
	require.NoError(t, swsc1.SetScalerParam1(6))
	require.Equal(t, float64(6), swsc1.ScalerParam1())

	require.NoError(t, swsc1.SetAlphaBlend(SoftwareScaleContextAlphaBlendNone))
	require.Equal(t, SoftwareScaleContextAlphaBlendNone, swsc1.AlphaBlend())
	require.NoError(t, swsc1.SetDestinationRange(1))
	require.Equal(t, 1, swsc1.DestinationRange())
	require.NoError(t, swsc1.SetDither(SoftwareScaleContextDitherBayer))
	require.Equal(t, SoftwareScaleContextDitherBayer, swsc1.Dither())
	require.NoError(t, swsc1.SetGammaFlag(1))
	require.Equal(t, 1, swsc1.GammaFlag())
	require.NoError(t, swsc1.SetIntent(SoftwareScaleContextIntentRelativeColorimetric))
	require.Equal(t, SoftwareScaleContextIntentRelativeColorimetric, swsc1.Intent())
	require.NoError(t, swsc1.SetOpaque(nil))
	require.True(t, swsc1.Opaque() == nil)
	require.NoError(t, swsc1.SetSourceRange(1))
	require.Equal(t, 1, swsc1.SourceRange())
	require.NoError(t, swsc1.SetThreads(4))
	require.Equal(t, 4, swsc1.Threads())

	sf := NewSoftwareScaleFilter(0, 0, 0, 0, 0, 0, 0)
	require.NotNil(t, sf)
	defer sf.Free()
	swsc3, err := CreateSoftwareScaleContext(srcW, srcH, srcPixelFormat, dstW, dstH, dstPixelFormat, swscf1, sf, nil, nil)
	require.NoError(t, err)
	defer swsc3.Free()

	f4, err := globalHelper.inputLastFrame("image-rgba.png", MediaTypeVideo, nil)
	require.NoError(t, err)

	f5 := AllocFrame()
	require.NotNil(t, f5)
	defer f5.Free()

	swsc2, err := CreateSoftwareScaleContext(f4.Width(), f4.Height(), f4.PixelFormat(), 512, 512, f4.PixelFormat(), NewSoftwareScaleContextFlags(SoftwareScaleContextFlagBilinear), nil, nil, nil)
	require.NoError(t, err)
	require.NoError(t, swsc2.ScaleFrame(f4, f5))

	b1, err := f5.Data().Bytes(1)
	require.NoError(t, err)

	b2, err := os.ReadFile("testdata/image-rgba-upscaled-bytes")
	require.NoError(t, err)
	require.Equal(t, b2, b1)
}
