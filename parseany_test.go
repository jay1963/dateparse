package dateparse

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Lets test to see how this performs using different Timezones/Locations
// Also of note, try changing your server/machine timezones and repeat
//
// !!!!! The time-zone of local machine effects the results!
// https://play.golang.org/p/IDHRalIyXh
// https://github.com/golang/go/issues/18012
func TestInLocation(t *testing.T) {

	denverLoc, err := time.LoadLocation("America/Denver")
	assert.Equal(t, nil, err)

	// Start out with time.UTC
	time.Local = time.UTC

	// Just normal parse to test out zone/offset
	ts := MustParse("2013-02-01 00:00:00")
	zone, offset := ts.Zone()
	assert.Equal(t, 0, offset, "Should have found offset = 0 %v", offset)
	assert.Equal(t, "UTC", zone, "Should have found zone = UTC %v", zone)
	assert.Equal(t, "2013-02-01 00:00:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	// Now lets set to denver (MST/MDT) and re-parse the same time string
	// and since no timezone info in string, we expect same result
	time.Local = denverLoc
	ts = MustParse("2013-02-01 00:00:00")
	zone, offset = ts.Zone()
	assert.Equal(t, 0, offset, "Should have found offset = 0 %v", offset)
	assert.Equal(t, "UTC", zone, "Should have found zone = UTC %v", zone)
	assert.Equal(t, "2013-02-01 00:00:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("Tue, 5 Jul 2017 16:28:13 -0700 (MST)")
	assert.Equal(t, "2017-07-05 23:28:13 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	// Now we are going to use ParseIn() and see that it gives different answer
	// with different zone, offset
	time.Local = nil
	ts, err = ParseIn("2013-02-01 00:00:00", denverLoc)
	assert.Equal(t, nil, err)
	zone, offset = ts.Zone()
	assert.Equal(t, -25200, offset, "Should have found offset = -25200 %v  %v", offset, denverLoc)
	assert.Equal(t, "MST", zone, "Should have found zone = MST %v", zone)
	assert.Equal(t, "2013-02-01 07:00:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts, err = ParseIn("18 January 2018", denverLoc)
	assert.Equal(t, nil, err)
	zone, offset = ts.Zone()
	assert.Equal(t, -25200, offset, "Should have found offset = 0 %v", offset)
	assert.Equal(t, "MST", zone, "Should have found zone = UTC %v", zone)
	assert.Equal(t, "2018-01-18 07:00:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	// Now we are going to use ParseLocal() and see that it gives same
	// answer as ParseIn when we have time.Local set to a location
	time.Local = denverLoc
	ts, err = ParseLocal("2013-02-01 00:00:00")
	assert.Equal(t, nil, err)
	zone, offset = ts.Zone()
	assert.Equal(t, -25200, offset, "Should have found offset = -25200 %v  %v", offset, denverLoc)
	assert.Equal(t, "MST", zone, "Should have found zone = MST %v", zone)
	assert.Equal(t, "2013-02-01 07:00:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	// Lets advance past daylight savings time start
	// use parseIn and see offset/zone has changed to Daylight Savings Equivalents
	ts, err = ParseIn("2013-04-01 00:00:00", denverLoc)
	assert.Equal(t, nil, err)
	zone, offset = ts.Zone()
	assert.Equal(t, -21600, offset, "Should have found offset = -21600 %v  %v", offset, denverLoc)
	assert.Equal(t, "MDT", zone, "Should have found zone = MDT %v", zone)
	assert.Equal(t, "2013-04-01 06:00:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	// reset to UTC
	time.Local = time.UTC

	//   UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
	ts = MustParse("Mon Jan  2 15:04:05 MST 2006")

	_, offset = ts.Zone()
	assert.Equal(t, 0, offset, "Should have found offset = 0 %v", offset)
	assert.Equal(t, "2006-01-02 15:04:05 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	// Now lets set to denver(mst/mdt)
	time.Local = denverLoc
	ts = MustParse("Mon Jan  2 15:04:05 MST 2006")

	// this time is different from one above parsed with time.Local set to UTC
	_, offset = ts.Zone()
	assert.Equal(t, -25200, offset, "Should have found offset = -25200 %v", offset)
	assert.Equal(t, "2006-01-02 22:04:05 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	// Now Reset To UTC
	time.Local = time.UTC

	// RFC850    = "Monday, 02-Jan-06 15:04:05 MST"
	ts = MustParse("Monday, 02-Jan-06 15:04:05 MST")
	_, offset = ts.Zone()
	assert.Equal(t, 0, offset, "Should have found offset = 0 %v", offset)
	assert.Equal(t, "2006-01-02 15:04:05 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	// Now lets set to denver
	time.Local = denverLoc
	ts = MustParse("Monday, 02-Jan-06 15:04:05 MST")
	_, offset = ts.Zone()
	assert.NotEqual(t, 0, offset, "Should have found offset %v", offset)
	assert.Equal(t, "2006-01-02 22:04:05 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))
}

func TestOne(t *testing.T) {
	time.Local = time.UTC
	var ts time.Time
	ts = MustParse("2014-05-11 08:20:13,787")
	assert.Equal(t, "2014-05-11 08:20:13.787 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))
}

type dateTest struct {
	in, out string
}

// {in: , out: },

var testInputs = []dateTest{
	{in: "oct 7, 1970", out: "1970-10-07 00:00:00 +0000 UTC"},
	{in: "oct 7, '70", out: "1970-10-07 00:00:00 +0000 UTC"},
	{in: "Oct 7, '70", out: "1970-10-07 00:00:00 +0000 UTC"},
	{in: "Feb 8, 2009 5:57:51 AM", out: "2009-02-08 05:57:51 +0000 UTC"},
	{in: "May 8, 2009 5:57:51 PM", out: "2009-05-08 17:57:51 +0000 UTC"},
	{in: "May 8, 2009 5:57:1 PM", out: "2009-05-08 17:57:01 +0000 UTC"},
	{in: "May 8, 2009 5:7:51 PM", out: "2009-05-08 17:07:51 +0000 UTC"},
	{in: "7 oct 70", out: "1970-10-07 00:00:00 +0000 UTC"},
	{in: "7 oct 1970", out: "1970-10-07 00:00:00 +0000 UTC"},
	//   ANSIC       = "Mon Jan _2 15:04:05 2006"
	{in: "Mon Jan  2 15:04:05 2006", out: "2006-01-02 15:04:05 +0000 UTC"},
	{in: "Thu May 8 17:57:51 2009", out: "2009-05-08 17:57:51 +0000 UTC"},
	{in: "Thu May  8 17:57:51 2009", out: "2009-05-08 17:57:51 +0000 UTC"},
	// RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
	{in: "Mon Jan 02 15:04:05 -0700 2006", out: "2006-01-02 22:04:05 +0000 UTC"},
	{in: "Thu May 08 17:57:51 -0700 2009", out: "2009-05-08 17:57:51 +0000 UTC"},
	//   UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
	{in: "Mon Jan  2 15:04:05 MST 2006", out: "2006-01-02 15:04:05 +0000 UTC"},
	{in: "Thu May  8 17:57:51 MST 2009", out: "2009-05-08 17:57:51 +0000 UTC"},
}

func TestParse(t *testing.T) {

	// Lets ensure we are operating on UTC
	time.Local = time.UTC

	zeroTime := time.Time{}.Unix()
	ts, err := ParseAny("INVALID")
	assert.Equal(t, zeroTime, ts.Unix())
	assert.NotEqual(t, nil, err)

	assert.Equal(t, true, testDidPanic("NOT GONNA HAPPEN"))

	for _, th := range testInputs {
		ts = MustParse(th.in)
		got := fmt.Sprintf("%v", ts.In(time.UTC))
		assert.Equal(t, th.out, got, "Expected %q but got %q from %q", th.out, got, th.in)
	}

	// RFC850    = "Monday, 02-Jan-06 15:04:05 MST"
	ts = MustParse("Wednesday, 07-May-09 08:00:43 MST")
	assert.Equal(t, "2009-05-07 08:00:43 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))
	ts = MustParse("Wednesday, 28-Feb-18 09:01:00 MST")
	assert.Equal(t, "2018-02-28 09:01:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	// ST_WEEKDAYCOMMADELTA
	//   Monday, 02 Jan 2006 15:04:05 -0700
	//   Monday, 02 Jan 2006 15:04:05 +0100
	ts = MustParse("Monday, 02 Jan 2006 15:04:05 +0100")
	assert.Equal(t, "2006-01-02 14:04:05 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))
	ts = MustParse("Monday, 02 Jan 2006 15:04:5 +0100")
	assert.Equal(t, "2006-01-02 14:04:05 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))
	ts = MustParse("Monday, 02 Jan 2006 15:4:05 +0100")
	assert.Equal(t, "2006-01-02 14:04:05 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))
	ts = MustParse("Monday, 02 Jan 2006 15:04:05 -0100")
	assert.Equal(t, "2006-01-02 16:04:05 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	// Another weird one, year on the end after UTC?
	ts = MustParse("Mon Aug 10 15:44:11 UTC+0000 2015")
	assert.Equal(t, "2015-08-10 15:44:11 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))
	ts = MustParse("Mon Aug 10 15:44:11 PST-0700 2015")
	assert.Equal(t, "2015-08-10 22:44:11 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))
	ts = MustParse("Mon Aug 10 15:44:11 CEST+0200 2015")
	assert.Equal(t, "2015-08-10 13:44:11 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	// Easily the worst Date format I have ever seen
	//  "Fri Jul 03 2015 18:04:07 GMT+0100 (GMT Daylight Time)"
	ts = MustParse("Fri Jul 03 2015 18:04:07 GMT+0100 (GMT Daylight Time)")
	assert.Equal(t, "2015-07-03 17:04:07 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("Fri, 03 Jul 2015 13:04:07 MST")
	assert.Equal(t, "2015-07-03 13:04:07 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("Mon, 2 Jan 2006 15:4:05 MST")
	assert.Equal(t, "2006-01-02 15:04:05 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("Mon, 2 Jan 2006 15:4:5 MST")
	assert.Equal(t, "2006-01-02 15:04:05 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("Mon, 02-Jan-06 15:04:05 MST")
	assert.Equal(t, "2006-01-02 15:04:05 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("Tue, 11 Jul 2017 16:28:13 +0200 (CEST)")
	assert.Equal(t, "2017-07-11 14:28:13 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("Tue, 5 Jul 2017 16:28:13 -0700 (CEST)")
	assert.Equal(t, "2017-07-05 23:28:13 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("Thu, 13 Jul 2017 08:58:40 +0100")
	assert.Equal(t, "2017-07-13 07:58:40 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("Mon, 02 Jan 2006 15:04:05 -0700")
	assert.Equal(t, "2006-01-02 22:04:05 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))
	ts = MustParse("Mon, 02 Jan 2006 15:4:05 -0700")
	assert.Equal(t, "2006-01-02 22:04:05 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))
	ts = MustParse("Mon, 02 Jan 2006 15:4:5 -0700")
	assert.Equal(t, "2006-01-02 22:04:05 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("Thu, 4 Jan 2018 17:53:36 +0000")
	assert.Equal(t, "2018-01-04 17:53:36 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	// not sure if this is anything close to a standard, never seen it before
	ts = MustParse("12 Feb 2006, 19:17")
	assert.Equal(t, "2006-02-12 19:17:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2 Feb 2006, 19:17")
	assert.Equal(t, "2006-02-02 19:17:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("12 Feb 2006, 19:17:22")
	assert.Equal(t, "2006-02-12 19:17:22 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2 Feb 2006, 19:17:22")
	assert.Equal(t, "2006-02-02 19:17:22 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("12 Feb 2006 19:17")
	assert.Equal(t, "2006-02-12 19:17:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2 Feb 2006 19:17")
	assert.Equal(t, "2006-02-02 19:17:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("12 Feb 2006 19:17:22")
	assert.Equal(t, "2006-02-12 19:17:22 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2 Feb 2006 19:17:22")
	assert.Equal(t, "2006-02-02 19:17:22 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	// // 19 Dec 2013 12:15:23 GMT
	// ts = MustParse("12 Feb 2006 19:17:22 GMT")
	// assert.Equal(t, "2006-02-12 19:17:22 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))
	// ts = MustParse("2 Feb 2006 19:17:22 GMT")
	// assert.Equal(t, "2006-02-02 19:17:22 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	// // 28 Mar 2010 15:45:30 +1100
	// ts = MustParse("12 Feb 2006 19:17:22")
	// assert.Equal(t, "2006-02-12 19:17:22 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2013-Feb-03")
	assert.Equal(t, "2013-02-03 00:00:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("03 February 2013")
	assert.Equal(t, "2013-02-03 00:00:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("3 February 2013")
	assert.Equal(t, "2013-02-03 00:00:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	//---------------------------------------------
	// Chinese 2014年04月18日

	ts = MustParse("2014年04月08日")
	assert.Equal(t, "2014-04-08 00:00:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2014年04月08日 19:17:22")
	assert.Equal(t, "2014-04-08 19:17:22 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	//---------------------------------------------
	//   mm/dd/yyyy ?

	ts = MustParse("3/31/2014")
	assert.Equal(t, "2014-03-31 00:00:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("03/31/2014")
	assert.Equal(t, "2014-03-31 00:00:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	// what type of date is this? 08/21/71
	ts = MustParse("08/21/71")
	assert.Equal(t, "1971-08-21 00:00:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	// m/d/yy
	ts = MustParse("8/1/71")
	assert.Equal(t, "1971-08-01 00:00:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("4/8/2014 22:05")
	assert.Equal(t, "2014-04-08 22:05:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))
	ts = MustParse("4/18/2014 22:05")
	assert.Equal(t, "2014-04-18 22:05:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))
	ts = MustParse("04/08/2014 22:05")
	assert.Equal(t, "2014-04-08 22:05:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("4/8/14 22:05")
	assert.Equal(t, "2014-04-08 22:05:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))
	ts = MustParse("4/18/14 22:05")
	assert.Equal(t, "2014-04-18 22:05:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))
	ts = MustParse("10/18/14 22:05")
	assert.Equal(t, "2014-10-18 22:05:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("04/2/2014 4:00:51")
	assert.Equal(t, "2014-04-02 04:00:51 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("8/8/1965 01:00:01 PM")
	assert.Equal(t, "1965-08-08 13:00:01 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("8/8/1965 12:00:01 AM")
	assert.Equal(t, "1965-08-08 00:00:01 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("8/8/1965 01:00 PM")
	assert.Equal(t, "1965-08-08 13:00:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("08/8/1965 01:00 PM")
	assert.Equal(t, "1965-08-08 13:00:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("8/08/1965 1:00 PM")
	assert.Equal(t, "1965-08-08 13:00:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("8/8/1965 12:00 AM")
	assert.Equal(t, "1965-08-08 00:00:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("8/13/1965 01:00 PM")
	assert.Equal(t, "1965-08-13 13:00:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("4/02/2014 03:00:51")
	assert.Equal(t, "2014-04-02 03:00:51 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("03/19/2012 10:11:59")
	assert.Equal(t, "2012-03-19 10:11:59 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("03/19/2012 10:11:59.3186369")
	assert.Equal(t, "2012-03-19 10:11:59.3186369 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	//---------------------------------------------
	//   yyyy/mm/dd ?

	ts = MustParse("2014/3/31")
	assert.Equal(t, "2014-03-31 00:00:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2014/03/31")
	assert.Equal(t, "2014-03-31 00:00:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2014/4/8 22:05")
	assert.Equal(t, "2014-04-08 22:05:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2014/4/8 2:05")
	assert.Equal(t, "2014-04-08 02:05:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2014/04/08 22:05")
	assert.Equal(t, "2014-04-08 22:05:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2014/04/2 03:00:51")
	assert.Equal(t, "2014-04-02 03:00:51 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2014/4/02 03:00:51")
	assert.Equal(t, "2014-04-02 03:00:51 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2012/03/19 10:11:59")
	assert.Equal(t, "2012-03-19 10:11:59 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2012/03/19 10:11:59.318")
	assert.Equal(t, "2012-03-19 10:11:59.318 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2012/03/19 10:11:59.3186369")
	assert.Equal(t, "2012-03-19 10:11:59.3186369 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2012/03/19 10:11:59.318636945")
	assert.Equal(t, "2012-03-19 10:11:59.318636945 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2012/03/19 10:11 PM")
	assert.Equal(t, "2012-03-19 22:11:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2012/03/19 1:11 PM")
	assert.Equal(t, "2012-03-19 13:11:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2012/3/19 10:11 PM")
	assert.Equal(t, "2012-03-19 22:11:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2012/3/3 10:11 PM")
	assert.Equal(t, "2012-03-03 22:11:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2012/03/19 10:11:59 PM")
	assert.Equal(t, "2012-03-19 22:11:59 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2012/3/3 10:11:59 PM")
	assert.Equal(t, "2012-03-03 22:11:59 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2012/03/03 10:11:59.345 PM")
	assert.Equal(t, "2012-03-03 22:11:59.345 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	//---------------------------------------------
	//   yyyy-mm-dd ?
	ts = MustParse("2009-08-12T22:15:09-07:00")
	assert.Equal(t, "2009-08-13 05:15:09 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2009-08-12T22:15:9-07:00")
	assert.Equal(t, "2009-08-13 05:15:09 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2009-08-12T22:15:09.123-07:00")
	assert.Equal(t, "2009-08-13 05:15:09.123 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)), "%v", ts.In(time.UTC))

	ts = MustParse("2009-08-12T22:15Z")
	assert.Equal(t, "2009-08-12 22:15:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)), "%v", ts.In(time.UTC))

	ts = MustParse("2009-08-12T22:15:09Z")
	assert.Equal(t, "2009-08-12 22:15:09 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)), "%v", ts.In(time.UTC))

	ts = MustParse("2009-08-12T22:15:09.99Z")
	assert.Equal(t, "2009-08-12 22:15:09.99 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2009-08-12T22:15:09.9999Z")
	assert.Equal(t, "2009-08-12 22:15:09.9999 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2009-08-12T22:15:09.99999999Z")
	assert.Equal(t, "2009-08-12 22:15:09.99999999 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2009-08-12T22:15:9.99999999Z")
	assert.Equal(t, "2009-08-12 22:15:09.99999999 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	// https://github.com/golang/go/issues/5294
	_, err = ParseAny(time.RFC3339)
	assert.NotEqual(t, nil, err)

	ts = MustParse("2009-08-12T22:15:09.123")
	assert.Equal(t, "2009-08-12 22:15:09.123 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2009-08-12T22:15:09.123456")
	assert.Equal(t, "2009-08-12 22:15:09.123456 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2009-08-12T22:15:09.12")
	assert.Equal(t, "2009-08-12 22:15:09.12 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2009-08-12T22:15:09.1")
	assert.Equal(t, "2009-08-12 22:15:09.1 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2009-08-12T22:15:09")
	assert.Equal(t, "2009-08-12 22:15:09 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2014-04-26 17:24:37.3186369")
	assert.Equal(t, "2014-04-26 17:24:37.3186369 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	//                  2015-06-25 01:25:37.115208593 +0000 UTC
	ts = MustParse("2012-08-03 18:31:59.257000000 +0000 UTC")
	assert.Equal(t, "2012-08-03 18:31:59.257 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2012-08-03 8:1:59.257000000 +0000 UTC")
	assert.Equal(t, "2012-08-03 08:01:59.257 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2012-8-03 18:31:59.257000000 +0000 UTC")
	assert.Equal(t, "2012-08-03 18:31:59.257 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2012-8-3 18:31:59.257000000 +0000 UTC")
	assert.Equal(t, "2012-08-03 18:31:59.257 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2015-09-30 18:48:56.35272715 +0000 UTC")
	assert.Equal(t, "2015-09-30 18:48:56.35272715 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2017-01-27 00:07:31.945167")
	assert.Equal(t, "2017-01-27 00:07:31.945167 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2012-08-03 18:31:59.257000000")
	assert.Equal(t, "2012-08-03 18:31:59.257 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2013-04-01 22:43")
	assert.Equal(t, "2013-04-01 22:43:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2013-04-01 22:43:22")
	assert.Equal(t, "2013-04-01 22:43:22 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2014-04-26 17:24:37.123")
	assert.Equal(t, "2014-04-26 17:24:37.123 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2014-04-26 17:24:37.123456 +0000 UTC")
	assert.Equal(t, "2014-04-26 17:24:37.123456 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2014-04-26 17:24:37.123456 UTC")
	assert.Equal(t, "2014-04-26 17:24:37.123456 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2014-04-26 17:24:37.123 UTC")
	assert.Equal(t, "2014-04-26 17:24:37.123 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2014-04-26 17:24:37.12 UTC")
	assert.Equal(t, "2014-04-26 17:24:37.12 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2014-04-26 17:24:37.1 UTC")
	assert.Equal(t, "2014-04-26 17:24:37.1 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2014-04-26 09:04:37.123 +0800")
	assert.Equal(t, "2014-04-26 01:04:37.123 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))
	ts = MustParse("2014-04-26 9:04:37.123 +0800")
	assert.Equal(t, "2014-04-26 01:04:37.123 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))
	ts = MustParse("2014-04-26 09:4:37.123 +0800")
	assert.Equal(t, "2014-04-26 01:04:37.123 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))
	ts = MustParse("2014-04-26 9:4:37.123 +0800")
	assert.Equal(t, "2014-04-26 01:04:37.123 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2014-04-26 17:24:37.123 -0800")
	assert.Equal(t, "2014-04-27 01:24:37.123 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2014-04-26 17:24:37.123456 +0800")
	assert.Equal(t, "2014-04-26 09:24:37.123456 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2014-04-26 17:24:37.123456 -0800")
	assert.Equal(t, "2014-04-27 01:24:37.123456 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2017-07-19 03:21:51+00:00")
	assert.Equal(t, "2017-07-19 03:21:51 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2017-07-09 03:01:51 +00:00 UTC")
	assert.Equal(t, "2017-07-09 03:01:51 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))
	ts = MustParse("2017-7-09 03:01:51 +00:00 UTC")
	assert.Equal(t, "2017-07-09 03:01:51 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))
	ts = MustParse("2017-07-9 03:01:51 +00:00 UTC")
	assert.Equal(t, "2017-07-09 03:01:51 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))
	ts = MustParse("2017-7-9 03:01:51 +00:00 UTC")
	assert.Equal(t, "2017-07-09 03:01:51 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2017-07-19 03:01:51 +00:00 UTC")
	assert.Equal(t, "2017-07-19 03:01:51 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))
	ts = MustParse("2017-07-19 3:01:51 +00:00 UTC")
	assert.Equal(t, "2017-07-19 03:01:51 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))
	ts = MustParse("2017-07-19 03:1:51 +00:00 UTC")
	assert.Equal(t, "2017-07-19 03:01:51 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))
	ts = MustParse("2017-07-19 3:1:51 +00:00 UTC")
	assert.Equal(t, "2017-07-19 03:01:51 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	// Golang Native Format
	ts = MustParse("2015-02-18 00:12:00 +0000 UTC")
	assert.Equal(t, "2015-02-18 00:12:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2015-02-18 00:12:00 +0000 GMT")
	assert.Equal(t, "2015-02-18 00:12:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2015-02-08 03:02:00 +0300 MSK")
	assert.Equal(t, "2015-02-08 00:02:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))
	ts = MustParse("2015-2-08 03:02:00 +0300 MSK")
	assert.Equal(t, "2015-02-08 00:02:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))
	ts = MustParse("2015-02-8 03:02:00 +0300 MSK")
	assert.Equal(t, "2015-02-08 00:02:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))
	ts = MustParse("2015-2-8 03:02:00 +0300 MSK")
	assert.Equal(t, "2015-02-08 00:02:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2014-12-16 06:20:00 UTC")
	assert.Equal(t, "2014-12-16 06:20:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2014-12-16 06:20:00 GMT")
	assert.Equal(t, "2014-12-16 06:20:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2014-12-16 06:20:00 +0000 UTC")
	assert.Equal(t, "2014-12-16 06:20:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2014-04-26 05:24:37 PM")
	assert.Equal(t, "2014-04-26 17:24:37 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	// This one is pretty special, it is TIMEZONE based but starts with P to emulate collions with PM
	ts = MustParse("2014-04-26 05:24:37 PST")
	assert.Equal(t, "2014-04-26 05:24:37 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2014-04-26")
	assert.Equal(t, "2014-04-26 00:00:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2014-04")
	assert.Equal(t, "2014-04-01 00:00:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2014-05-11 08:20:13,787")
	assert.Equal(t, "2014-05-11 08:20:13.787 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	_, err = ParseAny("2014-13-13 08:20:13,787") // month 13 doesn't exist so error
	assert.NotEqual(t, nil, err)

	ts = MustParse("2014-05-01 08:02:13 +00:00")
	assert.Equal(t, "2014-05-01 08:02:13 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))
	ts = MustParse("2014-5-01 08:02:13 +00:00")
	assert.Equal(t, "2014-05-01 08:02:13 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))
	ts = MustParse("2014-05-1 08:02:13 +00:00")
	assert.Equal(t, "2014-05-01 08:02:13 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))
	ts = MustParse("2014-5-1 08:02:13 +00:00")
	assert.Equal(t, "2014-05-01 08:02:13 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2014-05-11 08:20:13 +0000")
	assert.Equal(t, "2014-05-11 08:20:13 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2016-06-21T19:55:00+01:00")
	assert.Equal(t, "2016-06-21 18:55:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2016-06-21T19:55:00.799+01:00")
	assert.Equal(t, "2016-06-21 18:55:00.799 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2016-06-21T19:55:00+0100")
	assert.Equal(t, "2016-06-21 18:55:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2016-06-21T19:55:00-0700")
	assert.Equal(t, "2016-06-22 02:55:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("2016-06-21T19:55:00.799+0100")
	assert.Equal(t, "2016-06-21 18:55:00.799 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	//---------------------------------------------
	//   mm.dd.yyyy
	//   mm.dd.yy

	ts = MustParse("3.31.2014")
	assert.Equal(t, "2014-03-31 00:00:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("03.31.2014")
	assert.Equal(t, "2014-03-31 00:00:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("08.21.71")
	assert.Equal(t, "1971-08-21 00:00:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	//---------------------------------------------

	//  yyyymmdd and similar
	ts = MustParse("2014")
	assert.Equal(t, "2014-01-01 00:00:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("20140601")
	assert.Equal(t, "2014-06-01 00:00:00 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("1332151919")
	assert.Equal(t, "2012-03-19 10:11:59 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("1384216367111")
	assert.Equal(t, "2013-11-12 00:32:47.111 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts, _ = ParseIn("1384216367111", time.UTC)
	assert.Equal(t, "2013-11-12 00:32:47.111 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	ts = MustParse("1384216367111222")
	assert.Equal(t, "2013-11-12 00:32:47.111222 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	// Nanoseconds
	ts = MustParse("1384216367111222333")
	assert.Equal(t, "2013-11-12 00:32:47.111222333 +0000 UTC", fmt.Sprintf("%v", ts.In(time.UTC)))

	_, err = ParseAny("138421636711122233311111") // too many digits
	assert.NotEqual(t, nil, err)

	_, err = ParseAny("-1314")
	assert.NotEqual(t, nil, err)
}

func testDidPanic(datestr string) (paniced bool) {
	defer func() {
		if r := recover(); r != nil {
			paniced = true
		}
	}()
	MustParse(datestr)
	return false
}

func TestPStruct(t *testing.T) {

	denverLoc, err := time.LoadLocation("America/Denver")
	assert.Equal(t, nil, err)

	p := newParser("08.21.71", denverLoc)

	p.setMonth()
	assert.Equal(t, 0, p.moi)
	p.setDay()
	assert.Equal(t, 0, p.dayi)
	p.set(-1, "not")
	p.set(15, "not")
	assert.Equal(t, "08.21.71", p.datestr)
	assert.Equal(t, "08.21.71", string(p.format))
	assert.True(t, len(p.ds()) > 0)
	assert.True(t, len(p.ts()) > 0)
}
