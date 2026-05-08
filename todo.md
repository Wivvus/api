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
[x] Remove requirement for login after every redeploy
[x] Remove requirement to regrant location permission after every redeploy
[x] Allow attendees (but not the organiser), to rate an event out 5 after having attended it
[x] Remove distance from me filter
[x] All events to have multiple response options, defined by the creator (e.g. attend for 5k, attend for 10k, attend the post-run party)
    [x] Creator can add/remove named options during event creation and editing
    [x] First option is the default; removing an option remaps attendees to the first remaining option
    [x] Attendees section groups runners by their chosen option
[ ] the event list and map view page should poll the API every minute for new events
[ ] the event list should only retrieve the first 10 events, and paginate for more results
[ ] running groups - phase 1
    [ ] groups are created by a user
    [ ] groups have a name and a public information page
    [ ] groups have 1 owner
    [ ] groups have members
    [ ] members can apply to join, this requires approval from the owner. The owner can specify some questions to ask the new member when they apply to join.
    [ ] groups can have events
    [ ] members are notified by email if a new event is created in that group
    [ ] groups can have information that is only visible to group members
    [ ] In the account details page, add a tab to manage groups owned by this user
[ ] running groups - phase 2
    [ ] a group owner can assign other users as admins
    [ ] admins can create, edit and delete events in the group, and approve or deny member applications, but cannot modify or delete the group
    [ ] group members can rate a group (once per member, but editable) out of 5 with a comment.
    [ ] An event within a group can be repeating 
        [ ] events can repeat every day, week or month
        [ ] repeating events are created automatically 48 hours before it is scheduled
[ ] In the manage events section, split the view into 2 tabs:
    [ ] My events: the current view
    [ ] Attending events: all events in the past or future that this user is attending
[ ] In the account details page, add a tab to view a summary of my rating, and all of my individual ratings 
[x] Change URL scheme to wivvus.com/run/5 instead of wivvus.com/event/5
    [x] Old /events/:id URLs redirect to /run/:id
[x] allow run.wivvus.com/5 to point to wivvus.com/run/5