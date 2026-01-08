package astiav

//#include <libswscale/swscale.h>
import "C"
import (
	"errors"
	"unsafe"
)

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html
type SoftwareScaleContext struct {
	c *C.struct_SwsContext
}

// https://ffmpeg.org/doxygen/8.0/group__libsws.html#ga59cc19eff0434e7ec11676dc5e222ff3
func CreateSoftwareScaleContext(srcW, srcH int, srcFormat PixelFormat, dstW, dstH int, dstFormat PixelFormat, flags SoftwareScaleContextFlags, srcFilter, dstFilter *SoftwareScaleFilter, param []float64) (*SoftwareScaleContext, error) {
	ssc := &SoftwareScaleContext{}
	var cSrcFilter, cDstFilter *C.struct_SwsFilter
	if srcFilter != nil {
		cSrcFilter = srcFilter.c
	}
	if dstFilter != nil {
		cDstFilter = dstFilter.c
	}
	var cParam *C.double
	if len(param) > 0 {
		cParam = (*C.double)(unsafe.Pointer(&param[0]))
	}
	ssc.c = C.sws_getContext(
		C.int(srcW),
		C.int(srcH),
		C.enum_AVPixelFormat(srcFormat),
		C.int(dstW),
		C.int(dstH),
		C.enum_AVPixelFormat(dstFormat),
		C.int(flags),
		cSrcFilter,
		cDstFilter,
		cParam,
	)
	if ssc.c == nil {
		return nil, errors.New("astiav: empty new context")
	}

	classers.set(ssc)
	return ssc, nil
}

// https://ffmpeg.org/doxygen/8.0/structSwsFilter.html
type SoftwareScaleFilter struct {
	c *C.struct_SwsFilter
}

// https://ffmpeg.org/doxygen/8.0/group__libsws.html#ga70589e6382098677c77f022197609a56
func NewSoftwareScaleFilter(lumaGBlur, chromaGBlur, lumaSharpen, chromaSharpen, chromaHShift, chromaVShift float32, verbose int) *SoftwareScaleFilter {
	c := C.sws_getDefaultFilter(C.float(lumaGBlur), C.float(chromaGBlur), C.float(lumaSharpen), C.float(chromaSharpen), C.float(chromaHShift), C.float(chromaVShift), C.int(verbose))
	if c == nil {
		return nil
	}
	return &SoftwareScaleFilter{c: c}
}

// https://ffmpeg.org/doxygen/8.0/group__libsws.html#ga646067710384236814c1737f2a17688c
func (ssf *SoftwareScaleFilter) Free() {
	if ssf.c != nil {
		C.sws_freeFilter(ssf.c)
		ssf.c = nil
	}
}

// https://ffmpeg.org/doxygen/8.0/group__libsws.html#gad90b463ceeafdfd526994742f9954dbb
func (ssc *SoftwareScaleContext) Free() {
	if ssc.c != nil {
		// Make sure to clone the classer before freeing the object since
		// the C free method may reset the pointer
		c := newClonedClasser(ssc)
		C.sws_freeContext(ssc.c)
		ssc.c = nil
		// Make sure to remove from classers after freeing the object since
		// the C free method may use methods needing the classer
		if c != nil {
			classers.del(c)
		}

	}
}

var _ Classer = (*SoftwareScaleContext)(nil)

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a6866f52574bc730833d2580abc806261
func (ssc *SoftwareScaleContext) Class() *Class {
	if ssc.c == nil {
		return nil
	}
	return newClassFromC(unsafe.Pointer(ssc.c))
}

// https://ffmpeg.org/doxygen/8.0/swscale-v2_8txt.html#a20ffff3ac1378332422b93ed70264f4c
func (ssc *SoftwareScaleContext) ScaleFrame(src, dst *Frame) error {
	return newError(C.sws_scale_frame(ssc.c, dst.c, src.c))
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#aef45de443b59978fd38ad1531c618574
func (ssc *SoftwareScaleContext) Flags() SoftwareScaleContextFlags {
	if ssc.c == nil {
		return 0
	}
	return SoftwareScaleContextFlags(ssc.c.flags)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#aef45de443b59978fd38ad1531c618574
func (ssc *SoftwareScaleContext) SetFlags(swscf SoftwareScaleContextFlags) error {
	if ssc.c == nil {
		return errors.New("astiav: empty context")
	}
	ssc.c.flags = C.uint(swscf)
	return nil
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a883a891c8a2d4ea7a15a3a7055f64386
func (ssc *SoftwareScaleContext) DestinationWidth() int {
	if ssc.c == nil {
		return 0
	}
	return int(ssc.c.dst_w)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a883a891c8a2d4ea7a15a3a7055f64386
func (ssc *SoftwareScaleContext) SetDestinationWidth(i int) error {
	if ssc.c == nil {
		return errors.New("astiav: empty context")
	}
	ssc.c.dst_w = C.int(i)
	return nil
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a7facd34608c9258dae8c2942e3dce78f
func (ssc *SoftwareScaleContext) DestinationHeight() int {
	if ssc.c == nil {
		return 0
	}
	return int(ssc.c.dst_h)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a7facd34608c9258dae8c2942e3dce78f
func (ssc *SoftwareScaleContext) SetDestinationHeight(i int) error {
	if ssc.c == nil {
		return errors.New("astiav: empty context")
	}
	ssc.c.dst_h = C.int(i)
	return nil
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a0ff71c9ef5ab6dabf90378fa7bf836ec
func (ssc *SoftwareScaleContext) DestinationPixelFormat() PixelFormat {
	if ssc.c == nil {
		return 0
	}
	return PixelFormat(ssc.c.dst_format)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a0ff71c9ef5ab6dabf90378fa7bf836ec
func (ssc *SoftwareScaleContext) SetDestinationPixelFormat(p PixelFormat) error {
	if ssc.c == nil {
		return errors.New("astiav: empty context")
	}
	ssc.c.dst_format = C.int(p)
	return nil
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#aa7dc7a4f9ec57a7c37957259a51cd920
func (ssc *SoftwareScaleContext) SourceWidth() int {
	if ssc.c == nil {
		return 0
	}
	return int(ssc.c.src_w)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a0ff71c9ef5ab6dabf90378fa7bf836ec
func (ssc *SoftwareScaleContext) SetSourceWidth(i int) error {
	if ssc.c == nil {
		return errors.New("astiav: empty context")
	}
	ssc.c.src_w = C.int(i)
	return nil
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a0dbc8c02bd3b4cd472e07008009751ff
func (ssc *SoftwareScaleContext) SourceHeight() int {
	if ssc.c == nil {
		return 0
	}
	return int(ssc.c.src_h)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a0ff71c9ef5ab6dabf90378fa7bf836ec
func (ssc *SoftwareScaleContext) SetSourceHeight(i int) error {
	if ssc.c == nil {
		return errors.New("astiav: empty context")
	}
	ssc.c.src_h = C.int(i)
	return nil
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#aab113373f157ee3b255ad97481af0cd9
func (ssc *SoftwareScaleContext) SourcePixelFormat() PixelFormat {
	if ssc.c == nil {
		return 0
	}
	return PixelFormat(ssc.c.src_format)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#aab113373f157ee3b255ad97481af0cd9
func (ssc *SoftwareScaleContext) SetSourcePixelFormat(p PixelFormat) error {
	if ssc.c == nil {
		return errors.New("astiav: empty context")
	}
	ssc.c.src_format = C.int(p)
	return nil
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a90a448639c27486dc88b3ef4fa1252de
func (ssc *SoftwareScaleContext) Opaque() unsafe.Pointer {
	if ssc.c == nil {
		return nil
	}
	return ssc.c.opaque
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a0dbc8c02bd3b4cd472e07008009751ff
func (ssc *SoftwareScaleContext) SetOpaque(p unsafe.Pointer) error {
	if ssc.c == nil {
		return errors.New("astiav: empty context")
	}
	ssc.c.opaque = p
	return nil
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a3cc13e08b01c1152405b3a0e4313b255
func (ssc *SoftwareScaleContext) Threads() int {
	if ssc.c == nil {
		return 0
	}
	return int(ssc.c.threads)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a3cc13e08b01c1152405b3a0e4313b255
func (ssc *SoftwareScaleContext) SetThreads(i int) error {
	if ssc.c == nil {
		return errors.New("astiav: empty context")
	}
	ssc.c.threads = C.int(i)
	return nil
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#abdb353dd741ba1ff0bfab6bcb133c682
func (ssc *SoftwareScaleContext) Dither() SoftwareScaleContextDither {
	if ssc.c == nil {
		return 0
	}
	return SoftwareScaleContextDither(ssc.c.dither)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#abdb353dd741ba1ff0bfab6bcb133c682
func (ssc *SoftwareScaleContext) SetDither(d SoftwareScaleContextDither) error {
	if ssc.c == nil {
		return errors.New("astiav: empty context")
	}
	ssc.c.dither = C.SwsDither(d)
	return nil
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a8179d7c46f6e0acf4541abe42d03d743
func (ssc *SoftwareScaleContext) AlphaBlend() SoftwareScaleContextAlphaBlend {
	if ssc.c == nil {
		return 0
	}
	return SoftwareScaleContextAlphaBlend(ssc.c.alpha_blend)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a8179d7c46f6e0acf4541abe42d03d743
func (ssc *SoftwareScaleContext) SetAlphaBlend(a SoftwareScaleContextAlphaBlend) error {
	if ssc.c == nil {
		return errors.New("astiav: empty context")
	}
	ssc.c.alpha_blend = C.SwsAlphaBlend(a)
	return nil
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a28ec4f81a1e3dcce7c92f2e7be8e9bd1
func (ssc *SoftwareScaleContext) GammaFlag() int {
	if ssc.c == nil {
		return 0
	}
	return int(ssc.c.gamma_flag)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a28ec4f81a1e3dcce7c92f2e7be8e9bd1
func (ssc *SoftwareScaleContext) SetGammaFlag(i int) error {
	if ssc.c == nil {
		return errors.New("astiav: empty context")
	}
	ssc.c.gamma_flag = C.int(i)
	return nil
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#aebfe8c01f9ea0bb80c7ae4d433cc5062
func (ssc *SoftwareScaleContext) ScalerParam0() float64 {
	if ssc.c == nil {
		return 0
	}
	return float64(ssc.c.scaler_params[0])
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#aebfe8c01f9ea0bb80c7ae4d433cc5062
func (ssc *SoftwareScaleContext) SetScalerParam0(p float64) error {
	if ssc.c == nil {
		return errors.New("astiav: empty context")
	}
	ssc.c.scaler_params[0] = C.double(p)
	return nil
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#aebfe8c01f9ea0bb80c7ae4d433cc5062
func (ssc *SoftwareScaleContext) ScalerParam1() float64 {
	if ssc.c == nil {
		return 0
	}
	return float64(ssc.c.scaler_params[1])
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#aebfe8c01f9ea0bb80c7ae4d433cc5062
func (ssc *SoftwareScaleContext) SetScalerParam1(p float64) error {
	if ssc.c == nil {
		return errors.New("astiav: empty context")
	}
	ssc.c.scaler_params[1] = C.double(p)
	return nil
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#af10043377f39ca25a7af39f214c3fdff
func (ssc *SoftwareScaleContext) SourceRange() int {
	if ssc.c == nil {
		return 0
	}
	return int(ssc.c.src_range)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#af10043377f39ca25a7af39f214c3fdff
func (ssc *SoftwareScaleContext) SetSourceRange(i int) error {
	if ssc.c == nil {
		return errors.New("astiav: empty context")
	}
	ssc.c.src_range = C.int(i)
	return nil
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#ac8b4b166f254ccb8def6ce21deead82b
func (ssc *SoftwareScaleContext) DestinationRange() int {
	if ssc.c == nil {
		return 0
	}
	return int(ssc.c.dst_range)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#ac8b4b166f254ccb8def6ce21deead82b
func (ssc *SoftwareScaleContext) SetDestinationRange(i int) error {
	if ssc.c == nil {
		return errors.New("astiav: empty context")
	}
	ssc.c.dst_range = C.int(i)
	return nil
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#af1c9fa89159377a7949d9a1a9586f44f
func (ssc *SoftwareScaleContext) SourceHorizontalChromaPosition() int {
	if ssc.c == nil {
		return 0
	}
	return int(ssc.c.src_h_chr_pos)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#af1c9fa89159377a7949d9a1a9586f44f
func (ssc *SoftwareScaleContext) SetSourceHorizontalChromaPosition(i int) error {
	if ssc.c == nil {
		return errors.New("astiav: empty context")
	}
	ssc.c.src_h_chr_pos = C.int(i)
	return nil
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#af1c9fa89159377a7949d9a1a9586f44f
func (ssc *SoftwareScaleContext) SourceVerticalChromaPosition() int {
	if ssc.c == nil {
		return 0
	}
	return int(ssc.c.src_v_chr_pos)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#af1c9fa89159377a7949d9a1a9586f44f
func (ssc *SoftwareScaleContext) SetSourceVerticalChromaPosition(i int) error {
	if ssc.c == nil {
		return errors.New("astiav: empty context")
	}
	ssc.c.src_v_chr_pos = C.int(i)
	return nil
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a8d5fc2f5fa1e0f15b0a5b45b82b14c1a
func (ssc *SoftwareScaleContext) DestinationHorizontalChromaPosition() int {
	if ssc.c == nil {
		return 0
	}
	return int(ssc.c.dst_h_chr_pos)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a8d5fc2f5fa1e0f15b0a5b45b82b14c1a
func (ssc *SoftwareScaleContext) SetDestinationHorizontalChromaPosition(i int) error {
	if ssc.c == nil {
		return errors.New("astiav: empty context")
	}
	ssc.c.dst_h_chr_pos = C.int(i)
	return nil
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a8d5fc2f5fa1e0f15b0a5b45b82b14c1a
func (ssc *SoftwareScaleContext) DestinationVerticalChromaPosition() int {
	if ssc.c == nil {
		return 0
	}
	return int(ssc.c.dst_v_chr_pos)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a8d5fc2f5fa1e0f15b0a5b45b82b14c1a
func (ssc *SoftwareScaleContext) SetDestinationVerticalChromaPosition(i int) error {
	if ssc.c == nil {
		return errors.New("astiav: empty context")
	}
	ssc.c.dst_v_chr_pos = C.int(i)
	return nil
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a1c26a06608196ce7b73a34f7d20e2c13
func (ssc *SoftwareScaleContext) Intent() SoftwareScaleContextIntent {
	if ssc.c == nil {
		return 0
	}
	return SoftwareScaleContextIntent(ssc.c.intent)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a1c26a06608196ce7b73a34f7d20e2c13
func (ssc *SoftwareScaleContext) SetIntent(i SoftwareScaleContextIntent) error {
	if ssc.c == nil {
		return errors.New("astiav: empty context")
	}
	ssc.c.intent = C.int(i)
	return nil
}
