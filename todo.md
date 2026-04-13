[x] user registration (email/password with verification + Google OAuth account linking)
[x] Map view of local events (viewport-based fetching, refreshes on pan/zoom)
[x] Filters on list and map view
  [x] Pace (dual range slider, min/max)
  [x] length (dual range slider, min/max)
  [x] date (from date picker)
  [x] distance from me (single slider, Haversine radius with map overlay)
[x] allow an inclusive "slowest runner" pace option on events
[x] creator of event should be defaulted to attending
[x] delete account option
[x] Account overview and settings
[x] Update password
[x] Confirm event deletion
[x] logo and favicon
[x] set avatar for account
[x] custom 404 page
[x] add create event button to bottom of event list page
[x] metrics to see what users are most using and where they might be struggling with the site
[x] prompt to add the event to users calendar (e.g. google calendar or apple calendar) when choosing to attend an event
[x] see all the events I have created in a tab in my account (even ones in the past)
[x] send an email reminder to attendees, 12 hours before it is scheduled
[ ] the event list and map view page should poll the API every minute for new events
[x] Remove requirement for login after every redeploy
[x] Remove requirement to regrant location permission after every redeploy
[ ] Allow attendees (but not the organiser), to rate an event out 5 after having attended it
[x] Remove distance from me filter
[ ] Allow a user to create a running group
    [ ] groups have a name and a information page
    [ ] groups have 1 owner, but this owner can assign other users as admins
        [ ] admins can create, edit and delete events in the group, but cannot modify or delete the group
    [ ] groups have members
    [ ] members are notified if a new event is created in that group
    [ ] groups can have information that is only visible to group members
    [ ] group members can rate a group out of 5 with a comment.
    [ ] An event within a group can be repeating 
        [ ] events can repeat every day, week or month
        [ ] repeating events are created automatically 24 hours before it is scheduled
[ ] In the manage events section, split the view into 2 tabs:
    [ ] My events: the current view
    [ ] Attending events: all events in the past or future that this user is attending
[ ] In the account details page, add a tab to manage groups
[ ] In the account details page, add a tab to view a summary of my rating, and all of my individual ratings 